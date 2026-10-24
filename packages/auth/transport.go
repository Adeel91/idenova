package auth

import (
	"net/http"
)

type AuthenticatedTransport struct {
	Base     http.RoundTripper
	Provider *TokenProvider
}

func (t *AuthenticatedTransport) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	token, err := t.Provider.Token(request.Context())
	if err != nil {
		return nil, err
	}

	copy := request.Clone(request.Context())
	copy.Header.Set("Authorization", "Bearer "+token)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(copy)
}
