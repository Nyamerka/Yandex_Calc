package tests

import (
	"Yandex_Calc/routes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler_ValidExpressions(t *testing.T) {
	tests := []struct {
		expression string
		expected   string
	}{
		{"1+1", "2"},
		{"2*2", "4"},
		{"10/2", "5"},
		{"3-1", "2"},
	}

	router := routes.SetupRoutes()

	for _, tt := range tests {
		body := map[string]string{"expression": tt.expression}
		reqBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp routes.Response
		err := json.NewDecoder(rec.Body).Decode(&resp)
		if err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if resp.Result != tt.expected {
			t.Errorf("for expression %q, expected result %q, got %q", tt.expression, tt.expected, resp.Result)
		}
	}
}

func TestCalculateHandler_InvalidExpressions(t *testing.T) {
	tests := []string{
		"1//1",
		"abc",
		"",
	}

	router := routes.SetupRoutes()

	for _, expr := range tests {
		body := map[string]string{"expression": expr}
		reqBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d for expression %q", http.StatusUnprocessableEntity, rec.Code, expr)
		}

		var resp routes.Response
		err := json.NewDecoder(rec.Body).Decode(&resp)
		if err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if resp.Error != "Expression is not valid" {
			t.Errorf("expected error message 'Expression is not valid', got %q", resp.Error)
		}
	}
}

func TestCalculateHandler_InvalidMethod(t *testing.T) {
	router := routes.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
