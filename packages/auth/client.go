package auth

import (
	"errors"
	"net/http"
	"os"
	"time"
)

func NewServiceHTTPClient() (*http.Client, error) {
	secret := os.Getenv("IDENOVA_SERVICE_CLIENT_SECRET")
	if secret == "" {
		return nil, errors.New("IDENOVA_SERVICE_CLIENT_SECRET is required")
	}

	issuer := os.Getenv("KEYCLOAK_ISSUER")
	if issuer == "" {
		return nil, errors.New("KEYCLOAK_ISSUER is required")
	}

	provider, err := NewTokenProvider(
		issuer+"/protocol/openid-connect/token",
		"idenova-services",
		secret,
	)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &AuthenticatedTransport{
			Provider: provider,
		},
	}, nil
}
