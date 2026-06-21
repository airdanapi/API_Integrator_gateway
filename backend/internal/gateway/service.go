package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

// LogStore mendefinisikan kontrak penyimpanan log request gateway.
type LogStore interface {
	Insert(ctx context.Context, log model.RequestLog) (int64, error)
}

// Service mengimplementasikan logika bisnis gateway routing.
type Service struct {
	forwarder Forwarder
	logs      LogStore
	upstreams UpstreamConfig
	now       func() time.Time
}

// New membuat instance baru gateway Service.
func New(forwarder Forwarder, logs LogStore, upstreams UpstreamConfig, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{
		forwarder: forwarder,
		logs:      logs,
		upstreams: upstreams,
		now:       now,
	}
}

// ForwardPayment memproses POST /gateway/payment.
func (s *Service) ForwardPayment(ctx context.Context, claims auth.Claims, req PaymentRequest) (GatewayResponse, error) {
	if err := s.checkRole(claims); err != nil {
		return GatewayResponse{}, err
	}
	if err := req.Validate(); err != nil {
		return GatewayResponse{}, err
	}

	payload, _ := json.Marshal(req)
	result, err := s.forwarder.Forward(ctx, s.upstreams.SmartBankURL, payload)
	return s.buildResponse(ctx, claims, "/gateway/payment", payload, result, err, "payment")
}

// ForwardSmartBank memproses POST /gateway/smartbank.
func (s *Service) ForwardSmartBank(ctx context.Context, claims auth.Claims, req SmartBankRequest) (GatewayResponse, error) {
	if err := s.checkRole(claims); err != nil {
		return GatewayResponse{}, err
	}
	if err := req.Validate(); err != nil {
		return GatewayResponse{}, err
	}

	payload, _ := json.Marshal(req)
	result, err := s.forwarder.Forward(ctx, s.upstreams.SmartBankURL, payload)
	return s.buildResponse(ctx, claims, "/gateway/smartbank", payload, result, err, "smartbank")
}

// ForwardMarketplace memproses POST /gateway/marketplace.
func (s *Service) ForwardMarketplace(ctx context.Context, claims auth.Claims, req MarketplaceRequest) (GatewayResponse, error) {
	if err := s.checkRole(claims); err != nil {
		return GatewayResponse{}, err
	}
	if err := req.Validate(); err != nil {
		return GatewayResponse{}, err
	}

	payload, _ := json.Marshal(req)
	result, err := s.forwarder.Forward(ctx, s.upstreams.MarketplaceURL, payload)
	return s.buildResponse(ctx, claims, "/gateway/marketplace", payload, result, err, "marketplace")
}

// ForwardLogistics memproses POST /gateway/logistics.
func (s *Service) ForwardLogistics(ctx context.Context, claims auth.Claims, req LogisticsRequest) (GatewayResponse, error) {
	if err := s.checkRole(claims); err != nil {
		return GatewayResponse{}, err
	}
	if err := req.Validate(); err != nil {
		return GatewayResponse{}, err
	}

	payload, _ := json.Marshal(req)
	result, err := s.forwarder.Forward(ctx, s.upstreams.LogisticsURL, payload)
	return s.buildResponse(ctx, claims, "/gateway/logistics", payload, result, err, "logistics")
}

// ForwardSupplier memproses POST /gateway/supplier.
func (s *Service) ForwardSupplier(ctx context.Context, claims auth.Claims, req SupplierRequest) (GatewayResponse, error) {
	if err := s.checkRole(claims); err != nil {
		return GatewayResponse{}, err
	}
	if err := req.Validate(); err != nil {
		return GatewayResponse{}, err
	}

	payload, _ := json.Marshal(req)
	result, err := s.forwarder.Forward(ctx, s.upstreams.SupplierHubURL, payload)
	return s.buildResponse(ctx, claims, "/gateway/supplier", payload, result, err, "supplier")
}

func (s *Service) checkRole(claims auth.Claims) error {
	switch claims.Role {
	case model.RoleAdminGateway, model.RoleAppUser:
		return nil
	default:
		return ErrForbidden
	}
}

func (s *Service) buildResponse(
	ctx context.Context,
	claims auth.Claims,
	endpoint string,
	payload []byte,
	result ForwardResult,
	forwardErr error,
	idType string,
) (GatewayResponse, error) {
	now := s.now().UTC()
	transactionID := fmt.Sprintf("gw-%s-%d", idType, now.UnixNano())

	resp := GatewayResponse{
		TransactionID: transactionID,
		Message:       "request processed",
	}

	var httpStatus int
	var responseBody []byte

	if forwardErr != nil {
		if errors.Is(forwardErr, ErrUpstreamNotConfigured) {
			resp.Status = "success"
			resp.Message = "request accepted, upstream not configured"
			resp.Forwarded = false
			resp.Upstream = "not_configured"
			httpStatus = 200
		} else {
			resp.Status = "error"
			resp.Message = forwardErr.Error()
			resp.Forwarded = false
			resp.Upstream = "unreachable"
			httpStatus = 502
			s.logRequest(ctx, claims, endpoint, payload, httpStatus, responseBody, result.DurationMS, now)
			return GatewayResponse{}, forwardErr
		}
	} else {
		resp.Status = "success"
		resp.Message = "request forwarded successfully"
		resp.Forwarded = true
		resp.Upstream = s.upstreamName(endpoint)
		httpStatus = result.StatusCode
		responseBody = result.Body
		if result.Body != nil {
			var upstreamData interface{}
			if json.Unmarshal(result.Body, &upstreamData) == nil {
				resp.Data = upstreamData
			}
		}
	}

	durationMS := result.DurationMS
	s.logRequest(ctx, claims, endpoint, payload, httpStatus, responseBody, durationMS, now)

	return resp, nil
}

func (s *Service) upstreamName(endpoint string) string {
	switch endpoint {
	case "/gateway/payment", "/gateway/smartbank":
		return "smartbank"
	case "/gateway/marketplace":
		return "marketplace"
	case "/gateway/logistics":
		return "logistikit"
	case "/gateway/supplier":
		return "supplierhub"
	default:
		return "unknown"
	}
}

func (s *Service) logRequest(
	ctx context.Context,
	claims auth.Claims,
	endpoint string,
	payload []byte,
	status int,
	response []byte,
	durationMS int,
	ts time.Time,
) {
	if s.logs == nil {
		return
	}
	logEntry := model.RequestLog{
		Timestamp: ts,
		SourceApp: claims.AppName,
		Endpoint:  endpoint,
		Method:    "POST",
		Payload:   payload,
		Status:    status,
		Response:  response,
	}
	if durationMS > 0 {
		logEntry.DurationMS = &durationMS
	}
	s.logs.Insert(ctx, logEntry)
}
