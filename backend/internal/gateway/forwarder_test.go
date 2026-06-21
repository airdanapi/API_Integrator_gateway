package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPForwarder_EmptyURL(t *testing.T) {
	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	_, err := f.Forward(context.Background(), "", []byte(`{}`))
	if !errors.Is(err, ErrUpstreamNotConfigured) {
		t.Fatalf("expected ErrUpstreamNotConfigured, got %v", err)
	}
}

func TestHTTPForwarder_WhitespaceURL(t *testing.T) {
	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	_, err := f.Forward(context.Background(), "   ", []byte(`{}`))
	if !errors.Is(err, ErrUpstreamNotConfigured) {
		t.Fatalf("expected ErrUpstreamNotConfigured, got %v", err)
	}
}

func TestHTTPForwarder_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	result, err := f.Forward(context.Background(), server.URL, []byte(`{"test":true}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.StatusCode)
	}
	if string(result.Body) != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", result.Body)
	}
}

func TestHTTPForwarder_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal"}`))
	}))
	defer server.Close()

	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	result, err := f.Forward(context.Background(), server.URL, []byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", result.StatusCode)
	}
}

func TestHTTPForwarder_ConnectionRefusedRetry(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			// Close connection immediately to simulate connection error
			hijacker, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hijacker.Hijack()
				conn.Close()
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"recovered"}`))
	}))
	defer server.Close()

	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {} // no actual sleep

	result, err := f.Forward(context.Background(), server.URL, []byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.StatusCode)
	}
}

func TestHTTPForwarder_RetriesExhausted(t *testing.T) {
	// Use a closed server to simulate persistent connection failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := server.URL
	server.Close()

	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	_, err := f.Forward(context.Background(), closedURL, []byte(`{}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrUpstreamFailure) {
		t.Fatalf("expected ErrUpstreamFailure, got %v", err)
	}
}

func TestHTTPForwarder_DurationMeasurement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	currentTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	callCount := 0
	fakeNow := func() time.Time {
		callCount++
		if callCount == 1 {
			return currentTime
		}
		return currentTime.Add(150 * time.Millisecond)
	}

	f := NewHTTPForwarder(fakeNow)
	f.sleep = func(time.Duration) {}

	result, err := f.Forward(context.Background(), server.URL, []byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DurationMS != 150 {
		t.Fatalf("expected ~150ms duration, got %d", result.DurationMS)
	}
}

func TestHTTPForwarder_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	f := NewHTTPForwarder(time.Now)
	f.sleep = func(time.Duration) {}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := f.Forward(ctx, server.URL, []byte(`{}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
