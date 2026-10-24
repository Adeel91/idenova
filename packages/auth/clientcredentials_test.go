package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

func TestTokenProviderCachesToken(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)

			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			if r.Form.Get("grant_type") != "client_credentials" ||
				r.Form.Get("client_id") != "idenova-services" ||
				r.Form.Get("client_secret") != "test-secret" {
				t.Error("unexpected credentials request")
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "machine-token",
				"token_type":   "Bearer",
				"expires_in":   300,
			})
		},
	))
	defer server.Close()

	provider, err := NewTokenProvider(
		server.URL,
		"idenova-services",
		"test-secret",
	)
	if err != nil {
		t.Fatal(err)
	}

	const workers = 10
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			token, err := provider.Token(context.Background())
			if err != nil || token != "machine-token" {
				t.Errorf("unexpected token: %q, error: %v", token, err)
			}
		}()
	}
	wg.Wait()

	if got := requests.Load(); got != 1 {
		t.Fatalf("expected one token request, got %d", got)
	}
}

func TestTokenProviderRejectsFailedRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		},
	))
	defer server.Close()

	provider, err := NewTokenProvider(
		server.URL,
		"idenova-services",
		"wrong-secret",
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := provider.Token(context.Background()); err == nil {
		t.Fatal("expected token request failure")
	}
}
