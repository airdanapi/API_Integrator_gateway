package gateway

import (
	"errors"
	"strings"
)

// MaxPayloadSize adalah ukuran maksimum payload JSON yang diterima gateway (1 MB).
const MaxPayloadSize = 1 << 20

var (
	ErrInvalidPayload         = errors.New("invalid gateway payload")
	ErrUpstreamNotConfigured  = errors.New("upstream service not configured")
	ErrUpstreamFailure        = errors.New("upstream service failure")
	ErrForbidden              = errors.New("gateway access forbidden")
)

// PaymentRequest merepresentasikan payload untuk POST /gateway/payment.
type PaymentRequest struct {
	FromApp     string      `json:"from_app"`
	FromUser    string      `json:"from_user"`
	ToUser      string      `json:"to_user"`
	Amount      float64     `json:"amount"`
	Metadata    interface{} `json:"metadata"`
	ServiceType string      `json:"service_type"`
}

// Validate memvalidasi PaymentRequest sesuai kontrak PRD.
func (r PaymentRequest) Validate() error {
	if strings.TrimSpace(r.FromApp) == "" {
		return wrapValidation("from_app is required")
	}
	if strings.TrimSpace(r.FromUser) == "" {
		return wrapValidation("from_user is required")
	}
	if strings.TrimSpace(r.ToUser) == "" {
		return wrapValidation("to_user is required")
	}
	if r.Amount <= 0 {
		return wrapValidation("amount must be greater than zero")
	}
	if strings.TrimSpace(r.ServiceType) == "" {
		return wrapValidation("service_type is required")
	}
	return nil
}

// SmartBankRequest merepresentasikan payload untuk POST /gateway/smartbank.
type SmartBankRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Validate memvalidasi SmartBankRequest.
func (r SmartBankRequest) Validate() error {
	if strings.TrimSpace(r.Action) == "" {
		return wrapValidation("action is required")
	}
	return nil
}

// MarketplaceRequest merepresentasikan payload untuk POST /gateway/marketplace.
type MarketplaceRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Validate memvalidasi MarketplaceRequest.
func (r MarketplaceRequest) Validate() error {
	if strings.TrimSpace(r.Action) == "" {
		return wrapValidation("action is required")
	}
	return nil
}

// LogisticsRequest merepresentasikan payload untuk POST /gateway/logistics.
type LogisticsRequest struct {
	OrderID      string  `json:"order_id"`
	Address      string  `json:"address"`
	Distance     float64 `json:"distance"`
	ShippingType string  `json:"shipping_type"`
}

// Validate memvalidasi LogisticsRequest.
func (r LogisticsRequest) Validate() error {
	if strings.TrimSpace(r.OrderID) == "" {
		return wrapValidation("order_id is required")
	}
	if strings.TrimSpace(r.Address) == "" {
		return wrapValidation("address is required")
	}
	if r.Distance <= 0 {
		return wrapValidation("distance must be greater than zero")
	}
	if strings.TrimSpace(r.ShippingType) == "" {
		return wrapValidation("shipping_type is required")
	}
	return nil
}

// SupplierRequest merepresentasikan payload untuk POST /gateway/supplier.
type SupplierRequest struct {
	SupplierID string  `json:"supplier_id"`
	Material   string  `json:"material"`
	Qty        int     `json:"qty"`
	TotalCost  float64 `json:"total_cost"`
}

// Validate memvalidasi SupplierRequest.
func (r SupplierRequest) Validate() error {
	if strings.TrimSpace(r.SupplierID) == "" {
		return wrapValidation("supplier_id is required")
	}
	if strings.TrimSpace(r.Material) == "" {
		return wrapValidation("material is required")
	}
	if r.Qty <= 0 {
		return wrapValidation("qty must be greater than zero")
	}
	if r.TotalCost <= 0 {
		return wrapValidation("total_cost must be greater than zero")
	}
	return nil
}

// GatewayResponse adalah response standar dari semua endpoint gateway.
type GatewayResponse struct {
	Status        string      `json:"status"`
	TransactionID string      `json:"transaction_id,omitempty"`
	DeliveryID    string      `json:"delivery_id,omitempty"`
	OrderID       string      `json:"order_id,omitempty"`
	Message       string      `json:"message"`
	Forwarded     bool        `json:"forwarded"`
	Upstream      string      `json:"upstream"`
	Data          interface{} `json:"data,omitempty"`
}

func wrapValidation(msg string) error {
	return errors.Join(ErrInvalidPayload, errors.New(msg))
}
