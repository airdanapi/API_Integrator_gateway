package server_test

import (
	"encoding/json"
	"net/http"
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

func TestLandingEndpointContractIsPublic(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"})
	request := httptest.NewRequest(http.MethodGet, "/landing", nil)

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("GET /landing returned an unexpected error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /landing status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
		Data   struct {
			ServiceOverview struct {
				Name        string `json:"name"`
				Tagline     string `json:"tagline"`
				Description string `json:"description"`
				Benefits    []struct {
					Title       string `json:"title"`
					Description string `json:"description"`
				} `json:"benefits"`
			} `json:"service_overview"`
			ApplicationRoles []struct {
				Application string `json:"application"`
				Role        string `json:"role"`
				Interaction string `json:"interaction"`
			} `json:"application_roles"`
			IntegrationFlow []struct {
				Step        int    `json:"step"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"integration_flow"`
			ContactInfo struct {
				RepositoryURL string `json:"repository_url"`
				LoginStatus   string `json:"login_status"`
			} `json:"contact_info"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode landing response: %v", err)
	}

	if body.Status != "success" {
		t.Errorf("status = %q, want success", body.Status)
	}
	if body.Data.ServiceOverview.Name != "API Integrator Gateway" {
		t.Errorf("service name = %q, want API Integrator Gateway", body.Data.ServiceOverview.Name)
	}
	if body.Data.ServiceOverview.Tagline == "" || body.Data.ServiceOverview.Description == "" {
		t.Error("service overview tagline and description must not be empty")
	}
	if len(body.Data.ServiceOverview.Benefits) < 3 {
		t.Errorf("benefits count = %d, want at least 3", len(body.Data.ServiceOverview.Benefits))
	}

	wantApplications := []string{
		"SmartBank",
		"Marketplace",
		"POS",
		"SupplierHub",
		"LogistiKita",
		"UMKM Insight",
		"API Gateway",
	}
	if len(body.Data.ApplicationRoles) != len(wantApplications) {
		t.Fatalf(
			"application roles count = %d, want %d",
			len(body.Data.ApplicationRoles),
			len(wantApplications),
		)
	}
	for index, wantApplication := range wantApplications {
		role := body.Data.ApplicationRoles[index]
		if role.Application != wantApplication {
			t.Errorf(
				"application role %d = %q, want %q",
				index,
				role.Application,
				wantApplication,
			)
		}
		if role.Role == "" || role.Interaction == "" {
			t.Errorf("application role %q must include role and interaction", role.Application)
		}
	}

	if len(body.Data.IntegrationFlow) < 4 {
		t.Errorf("integration flow count = %d, want at least 4", len(body.Data.IntegrationFlow))
	}
	for index, flowStep := range body.Data.IntegrationFlow {
		if flowStep.Step != index+1 {
			t.Errorf("integration flow step %d = %d, want %d", index, flowStep.Step, index+1)
		}
		if flowStep.Title == "" || flowStep.Description == "" {
			t.Errorf("integration flow step %d must include title and description", index)
		}
	}

	if body.Data.ContactInfo.RepositoryURL != "https://github.com/airdanapi/API_Integrator_gateway" {
		t.Errorf("repository URL = %q, want official repository", body.Data.ContactInfo.RepositoryURL)
	}
	if body.Data.ContactInfo.LoginStatus != "coming_soon" {
		t.Errorf("login status = %q, want coming_soon", body.Data.ContactInfo.LoginStatus)
	}
}

func TestLandingEndpointAllowsCrossOriginRequests(t *testing.T) {
	app := server.NewApp(config.Config{AppEnv: "test"})
	request := httptest.NewRequest(http.MethodGet, "/landing", nil)
	request.Header.Set("Origin", "https://portal.example.test")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("GET /landing returned an unexpected error: %v", err)
	}
	defer response.Body.Close()

	if response.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(
			"Access-Control-Allow-Origin = %q, want *",
			response.Header.Get("Access-Control-Allow-Origin"),
		)
	}
}
