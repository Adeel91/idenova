package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testVerifier(t *testing.T) (*Verifier, *rsa.PrivateKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return &Verifier{
		issuer:   "http://localhost:8085/realms/idenova",
		audience: "idenova-api",
		key: func(token *jwt.Token) (any, error) {
			return &key.PublicKey, nil
		},
	}, key
}

func signTestToken(
	t *testing.T,
	key *rsa.PrivateKey,
	role, issuer, audience, authorizedParty string,
	expires time.Time,
) string {
	t.Helper()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   "test-user",
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		AuthorizedParty: authorizedParty,
	}
	claims.RealmAccess.Roles = []string{role}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

func TestTokenValidation(t *testing.T) {
	v, key := testVerifier(t)

	tests := []struct {
		name, role, issuer, audience, party string
		expires                             time.Time
		valid                               bool
	}{
		{"valid", "ADMIN", v.issuer, v.audience, "idenova-dashboard", time.Now().Add(time.Hour), true},
		{"expired", "ADMIN", v.issuer, v.audience, "idenova-dashboard", time.Now().Add(-time.Hour), false},
		{"wrong issuer", "ADMIN", "http://other", v.audience, "idenova-dashboard", time.Now().Add(time.Hour), false},
		{"wrong audience", "ADMIN", v.issuer, "other-api", "idenova-dashboard", time.Now().Add(time.Hour), false},
		{"wrong client", "ADMIN", v.issuer, v.audience, "other-client", time.Now().Add(time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signed := signTestToken(
				t, key, tt.role, tt.issuer, tt.audience,
				tt.party, tt.expires,
			)

			_, err := v.Verify(signed)
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v, verification error=%v", tt.valid, err)
			}
		})
	}
}

func TestAuthorizationMiddleware(t *testing.T) {
	v, key := testVerifier(t)

	handler := v.RequireRoles("ADMIN")(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok || claims.Subject != "test-user" {
				t.Error("missing authenticated claims")
			}
			w.WriteHeader(http.StatusNoContent)
		},
	))

	tests := []struct {
		name, token string
		status      int
	}{
		{"missing token", "", http.StatusUnauthorized},
		{"malformed token", "not-a-token", http.StatusUnauthorized},
		{
			"reviewer forbidden",
			signTestToken(t, key, "REVIEWER", v.issuer, v.audience,
				"idenova-dashboard", time.Now().Add(time.Hour)),
			http.StatusForbidden,
		},
		{
			"admin allowed",
			signTestToken(t, key, "ADMIN", v.issuer, v.audience,
				"idenova-dashboard", time.Now().Add(time.Hour)),
			http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("expected %d got %d", tt.status, rec.Code)
			}
		})
	}

	_ = strings.TrimSpace
}
