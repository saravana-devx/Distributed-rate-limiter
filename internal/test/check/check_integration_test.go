package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ratelimiter/internal/bootstrap"
	"ratelimiter/internal/config"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	cfg := &config.Config{
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUser:     "postgres",
		PostgresPassword: "process",
		PostgresDB:       "ratelimiter",
		RedisAddr:        "localhost:6379",
		RedisPassword:    "ratelimiter",
		ServerPort:       "8080",
	}
	app, err := bootstrap.New(cfg)
	if err != nil {
		t.Fatalf("bootstrap.New: %v", err)
	}
	return app.Router()
}

type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func doRequest(t *testing.T, router http.Handler, method, path string, body any, headers map[string]string) (*httptest.ResponseRecorder, envelope) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var env envelope
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
	return rec, env
}

func TestCheckEndpoint_HitsLimit(t *testing.T) {
	router := newTestRouter(t)

	clientID := fmt.Sprintf("test-check-%d", time.Now().UnixNano())

	createReq := map[string]any{
		"clientId":       clientID,
		"algorithm":      "fixed_window",
		"limit":          2,
		"window_seconds": 60,
	}
	rec, env := doRequest(t, router, http.MethodPost, "/api/v1/clients", createReq, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create client: expected 201, got %d (%s)", rec.Code, rec.Body.String())
	}

	var created struct {
		ClientID string
		APIKey   string
	}
	if err := json.Unmarshal(env.Data, &created); err != nil {
		t.Fatalf("decode created client: %v", err)
	}
	defer doRequest(t, router, http.MethodDelete, "/api/v1/clients/"+clientID, nil, nil)

	checkBody := map[string]any{"identifier": "user-1"}
	headers := map[string]string{"X-API-Key": created.APIKey}

	for i := range 2 {
		rec, _ := doRequest(t, router, http.MethodPost, "/api/v1/check", checkBody, headers)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d (%s)", i+1, rec.Code, rec.Body.String())
		}
	}

	rec, _ = doRequest(t, router, http.MethodPost, "/api/v1/check", checkBody, headers)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after limit exhausted, got %d (%s)", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Errorf("expected Retry-After header on 429 response")
	}
}
