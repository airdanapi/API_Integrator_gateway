package gateway

import (
	"errors"
	"testing"
)

func TestPaymentRequestValidate(t *testing.T) {
	valid := PaymentRequest{
		FromApp:     "Marketplace",
		FromUser:    "user1",
		ToUser:      "user2",
		Amount:      10000,
		ServiceType: "payment",
	}

	t.Run("valid request", func(t *testing.T) {
		if err := valid.Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	tests := []struct {
		name string
		mod  func(PaymentRequest) PaymentRequest
	}{
		{"missing from_app", func(r PaymentRequest) PaymentRequest { r.FromApp = ""; return r }},
		{"missing from_user", func(r PaymentRequest) PaymentRequest { r.FromUser = ""; return r }},
		{"missing to_user", func(r PaymentRequest) PaymentRequest { r.ToUser = ""; return r }},
		{"zero amount", func(r PaymentRequest) PaymentRequest { r.Amount = 0; return r }},
		{"negative amount", func(r PaymentRequest) PaymentRequest { r.Amount = -1; return r }},
		{"missing service_type", func(r PaymentRequest) PaymentRequest { r.ServiceType = ""; return r }},
		{"whitespace from_app", func(r PaymentRequest) PaymentRequest { r.FromApp = "  "; return r }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mod(valid).Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidPayload) {
				t.Fatalf("expected ErrInvalidPayload, got %v", err)
			}
		})
	}
}

func TestSmartBankRequestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r := SmartBankRequest{Action: "check_balance"}
		if err := r.Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
	t.Run("missing action", func(t *testing.T) {
		r := SmartBankRequest{}
		if err := r.Validate(); !errors.Is(err, ErrInvalidPayload) {
			t.Fatalf("expected ErrInvalidPayload, got %v", err)
		}
	})
	t.Run("whitespace action", func(t *testing.T) {
		r := SmartBankRequest{Action: "   "}
		if err := r.Validate(); !errors.Is(err, ErrInvalidPayload) {
			t.Fatalf("expected ErrInvalidPayload, got %v", err)
		}
	})
}

func TestMarketplaceRequestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r := MarketplaceRequest{Action: "get_order"}
		if err := r.Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
	t.Run("missing action", func(t *testing.T) {
		r := MarketplaceRequest{}
		if err := r.Validate(); !errors.Is(err, ErrInvalidPayload) {
			t.Fatalf("expected ErrInvalidPayload, got %v", err)
		}
	})
}

func TestLogisticsRequestValidate(t *testing.T) {
	valid := LogisticsRequest{
		OrderID:      "ORD-001",
		Address:      "Jl. Merdeka No. 1",
		Distance:     5.5,
		ShippingType: "standard",
	}

	t.Run("valid", func(t *testing.T) {
		if err := valid.Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	tests := []struct {
		name string
		mod  func(LogisticsRequest) LogisticsRequest
	}{
		{"missing order_id", func(r LogisticsRequest) LogisticsRequest { r.OrderID = ""; return r }},
		{"missing address", func(r LogisticsRequest) LogisticsRequest { r.Address = ""; return r }},
		{"zero distance", func(r LogisticsRequest) LogisticsRequest { r.Distance = 0; return r }},
		{"negative distance", func(r LogisticsRequest) LogisticsRequest { r.Distance = -1; return r }},
		{"missing shipping_type", func(r LogisticsRequest) LogisticsRequest { r.ShippingType = ""; return r }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mod(valid).Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidPayload) {
				t.Fatalf("expected ErrInvalidPayload, got %v", err)
			}
		})
	}
}

func TestSupplierRequestValidate(t *testing.T) {
	valid := SupplierRequest{
		SupplierID: "SUP-001",
		Material:   "Beras",
		Qty:        100,
		TotalCost:  500000,
	}

	t.Run("valid", func(t *testing.T) {
		if err := valid.Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	tests := []struct {
		name string
		mod  func(SupplierRequest) SupplierRequest
	}{
		{"missing supplier_id", func(r SupplierRequest) SupplierRequest { r.SupplierID = ""; return r }},
		{"missing material", func(r SupplierRequest) SupplierRequest { r.Material = ""; return r }},
		{"zero qty", func(r SupplierRequest) SupplierRequest { r.Qty = 0; return r }},
		{"negative qty", func(r SupplierRequest) SupplierRequest { r.Qty = -5; return r }},
		{"zero total_cost", func(r SupplierRequest) SupplierRequest { r.TotalCost = 0; return r }},
		{"negative total_cost", func(r SupplierRequest) SupplierRequest { r.TotalCost = -1; return r }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mod(valid).Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidPayload) {
				t.Fatalf("expected ErrInvalidPayload, got %v", err)
			}
		})
	}
}
