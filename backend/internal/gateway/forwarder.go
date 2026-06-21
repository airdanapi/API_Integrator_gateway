package gateway

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	upstreamTimeout   = 10 * time.Second
	maxRetries        = 2
	baseRetryInterval = 100 * time.Millisecond
)

// UpstreamConfig menyimpan URL layanan upstream.
// URL yang kosong berarti layanan belum dikonfigurasi.
type UpstreamConfig struct {
	SmartBankURL   string
	MarketplaceURL string
	LogisticsURL   string
	SupplierHubURL string
}

// ForwardResult menyimpan hasil forwarding ke upstream.
type ForwardResult struct {
	StatusCode int
	Body       []byte
	DurationMS int
}

// Forwarder mendefinisikan kontrak untuk meneruskan request ke upstream.
type Forwarder interface {
	Forward(ctx context.Context, upstreamURL string, payload []byte) (ForwardResult, error)
}

// HTTPForwarder mengimplementasikan Forwarder menggunakan HTTP client.
type HTTPForwarder struct {
	client *http.Client
	now    func() time.Time
	sleep  func(time.Duration)
}

// NewHTTPForwarder membuat instance baru HTTPForwarder.
func NewHTTPForwarder(now func() time.Time) *HTTPForwarder {
	if now == nil {
		now = time.Now
	}
	return &HTTPForwarder{
		client: &http.Client{Timeout: upstreamTimeout},
		now:    now,
		sleep:  time.Sleep,
	}
}

// Forward mengirim JSON payload ke upstream URL dengan retry logic.
func (f *HTTPForwarder) Forward(ctx context.Context, upstreamURL string, payload []byte) (ForwardResult, error) {
	upstreamURL = strings.TrimSpace(upstreamURL)
	if upstreamURL == "" {
		return ForwardResult{}, ErrUpstreamNotConfigured
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			interval := baseRetryInterval * time.Duration(1<<uint(attempt-1))
			if f.sleep != nil {
				f.sleep(interval)
			}
		}

		result, err := f.doForward(ctx, upstreamURL, payload)
		if err == nil {
			return result, nil
		}

		if !isRetryable(err) {
			return ForwardResult{}, fmt.Errorf("%w: %v", ErrUpstreamFailure, err)
		}
		lastErr = err
	}
	return ForwardResult{}, fmt.Errorf("%w: retries exhausted: %v", ErrUpstreamFailure, lastErr)
}

func (f *HTTPForwarder) doForward(ctx context.Context, url string, payload []byte) (ForwardResult, error) {
	start := f.now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return ForwardResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return ForwardResult{}, err
	}
	defer resp.Body.Close()

	var body bytes.Buffer
	if _, err := body.ReadFrom(resp.Body); err != nil {
		return ForwardResult{}, err
	}

	elapsed := f.now().Sub(start)
	return ForwardResult{
		StatusCode: resp.StatusCode,
		Body:       body.Bytes(),
		DurationMS: int(elapsed.Milliseconds()),
	}, nil
}

// isRetryable memeriksa apakah error bersifat transient (koneksi/timeout).
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	// connection refused, DNS error, EOF, etc.
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "EOF")
}
