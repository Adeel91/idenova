package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticatedTransport(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "test-service-token",
				"token_type":   "Bearer",
				"expires_in":   300,
			})
		},
	))
	defer tokenServer.Close()

	serviceServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer test-service-token" {
				t.Error("missing service authorization")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		},
	))
	defer serviceServer.Close()

	provider, err := NewTokenProvider(
		tokenServer.URL,
		"idenova-services",
		"test-secret",
	)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		Transport: &AuthenticatedTransport{
			Provider: provider,
		},
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		serviceServer.URL,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.StatusCode)
	}
}
