package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/auth"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/server"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type stubLoginService struct {
	result auth.LoginResult
	err    error
}

func (stub stubLoginService) Login(
	context.Context,
	auth.LoginRequest,
) (auth.LoginResult, error) {
	return stub.result, stub.err
}

type stubTokenVerifier struct {
	claims auth.Claims
	err    error
	token  string
}

func (stub *stubTokenVerifier) Validate(token string) (auth.Claims, error) {
	stub.token = token
	return stub.claims, stub.err
}

func TestLoginEndpointReturnsTokenContract(t *testing.T) {
	app := server.NewApp(testConfig(), server.Dependencies{
		AuthService: stubLoginService{result: auth.LoginResult{
			Token:        "signed-token",
			Role:         model.RoleAdminGateway,
			AppName:      "API Gateway",
			DashboardURL: "/dashboard/admin",
			ExpiresIn:    3600,
		}},
	})

	response := performJSONRequest(t, app, http.MethodPost, "/auth/login", map[string]string{
		"username": "admin",
		"password": "admin-password",
		"app_name": "API Gateway",
	}, "")
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("POST /auth/login status = %d, want 200", response.StatusCode)
	}

	var body struct {
		Status string           `json:"status"`
		Data   auth.LoginResult `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if body.Status != "success" || body.Data.Token != "signed-token" {
		t.Fatalf("unexpected login response: %#v", body)
	}
}

func TestLoginEndpointValidatesPayloadAndUsesGenericUnauthorizedResponse(t *testing.T) {
	app := server.NewApp(testConfig(), server.Dependencies{
		AuthService: stubLoginService{err: auth.ErrInvalidCredentials},
	})

	missingField := performJSONRequest(
		t,
		app,
		http.MethodPost,
		"/auth/login",
		map[string]string{"username": "admin"},
		"",
	)
	defer missingField.Body.Close()
	if missingField.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing field status = %d, want 400", missingField.StatusCode)
	}

	invalidCredentials := performJSONRequest(t, app, http.MethodPost, "/auth/login", map[string]string{
		"username": "admin",
		"password": "wrong",
		"app_name": "API Gateway",
	}, "")
	defer invalidCredentials.Body.Close()
	if invalidCredentials.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid credentials status = %d, want 401", invalidCredentials.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(invalidCredentials.Body).Decode(&body); err != nil {
		t.Fatalf("decode unauthorized response: %v", err)
	}
	encoded, _ := json.Marshal(body)
	if bytes.Contains(bytes.ToLower(encoded), []byte("username")) ||
		bytes.Contains(bytes.ToLower(encoded), []byte("password")) {
		t.Fatalf("unauthorized response leaks credential detail: %s", encoded)
	}
}

func TestAuthMeRequiresBearerTokenAndReturnsClaims(t *testing.T) {
	verifier := &stubTokenVerifier{claims: auth.Claims{
		Username: "admin",
		Role:     model.RoleAdminGateway,
		AppName:  "API Gateway",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "42",
			Issuer:    "api-integrator-gateway",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}}
	app := server.NewApp(testConfig(), server.Dependencies{
		TokenVerifier: verifier,
	})

	withoutToken := performJSONRequest(t, app, http.MethodGet, "/auth/me", nil, "")
	defer withoutToken.Body.Close()
	if withoutToken.StatusCode != http.StatusUnauthorized {
		t.Fatalf("GET /auth/me without token status = %d, want 401", withoutToken.StatusCode)
	}

	withToken := performJSONRequest(
		t,
		app,
		http.MethodGet,
		"/auth/me",
		nil,
		"Bearer signed-token",
	)
	defer withToken.Body.Close()
	if withToken.StatusCode != http.StatusOK {
		t.Fatalf("GET /auth/me status = %d, want 200", withToken.StatusCode)
	}
	if verifier.token != "signed-token" {
		t.Fatalf("validated token = %q, want signed-token", verifier.token)
	}

	var body struct {
		Status string `json:"status"`
		Data   struct {
			UserID   string     `json:"user_id"`
			Username string     `json:"username"`
			Role     model.Role `json:"role"`
			AppName  string     `json:"app_name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(withToken.Body).Decode(&body); err != nil {
		t.Fatalf("decode /auth/me response: %v", err)
	}
	if body.Data.UserID != "42" ||
		body.Data.Username != "admin" ||
		body.Data.Role != model.RoleAdminGateway ||
		body.Data.AppName != "API Gateway" {
		t.Fatalf("unexpected /auth/me response: %#v", body)
	}
}

func TestAuthMeRejectsMalformedAndInvalidTokens(t *testing.T) {
	verifier := &stubTokenVerifier{err: errors.New("invalid token")}
	app := server.NewApp(testConfig(), server.Dependencies{TokenVerifier: verifier})

	for _, authorization := range []string{
		"signed-token",
		"Basic signed-token",
		"Bearer",
		"Bearer invalid-token",
	} {
		response := performJSONRequest(
			t,
			app,
			http.MethodGet,
			"/auth/me",
			nil,
			authorization,
		)
		defer response.Body.Close()
		if response.StatusCode != http.StatusUnauthorized {
			t.Errorf(
				"Authorization %q status = %d, want 401",
				authorization,
				response.StatusCode,
			)
		}
	}
}

func testConfig() config.Config {
	return config.Config{
		AppEnv:    "test",
		JWTIssuer: "api-integrator-gateway",
	}
}

func performJSONRequest(
	t *testing.T,
	app interface {
		Test(*http.Request, ...fiber.TestConfig) (*http.Response, error)
	},
	method string,
	path string,
	body any,
	authorization string,
) *http.Response {
	t.Helper()

	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
	}

	request := httptest.NewRequest(method, path, bytes.NewReader(requestBody))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("%s %s returned an unexpected error: %v", method, path, err)
	}
	return response
}
