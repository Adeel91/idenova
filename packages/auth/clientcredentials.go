package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TokenProvider struct {
	endpoint     string
	clientID     string
	clientSecret string
	client       *http.Client

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func NewTokenProvider(
	endpoint, clientID, clientSecret string,
) (*TokenProvider, error) {
	if endpoint == "" || clientID == "" || clientSecret == "" {
		return nil, errors.New("token endpoint, client ID and secret are required")
	}

	return &TokenProvider{
		endpoint:     endpoint,
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (p *TokenProvider) Token(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Before(p.expiresAt) {
		return p.token, nil
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request service token: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned HTTP %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode service token: %w", err)
	}

	if result.AccessToken == "" ||
		!strings.EqualFold(result.TokenType, "Bearer") ||
		result.ExpiresIn <= 0 {
		return "", errors.New("invalid service token response")
	}

	lifetime := time.Duration(result.ExpiresIn) * time.Second
	margin := min(30*time.Second, lifetime/10)

	p.token = result.AccessToken
	p.expiresAt = time.Now().Add(lifetime - margin)

	return p.token, nil
}
