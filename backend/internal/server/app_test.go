package server_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
)

func TestHealthEndpointContract(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"})

	response, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatalf("GET /health returned an unexpected error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Fatalf("GET /health status = %d, want 200", response.StatusCode)
	}

	var body struct {
		Status string `json:"status"`
		Data   struct {
			Service     string `json:"service"`
			Environment string `json:"environment"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if body.Status != "success" {
		t.Errorf("status = %q, want success", body.Status)
	}
	if body.Data.Service != "api-integrator-gateway" {
		t.Errorf("service = %q, want api-integrator-gateway", body.Data.Service)
	}
	if body.Data.Environment != "test" {
		t.Errorf("environment = %q, want test", body.Data.Environment)
	}
}
