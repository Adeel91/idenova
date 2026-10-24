package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	AuthorizedParty string `json:"azp"`
	RealmAccess     struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func (c *Claims) HasRole(role string) bool {
	for _, r := range c.RealmAccess.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type Verifier struct {
	issuer   string
	audience string
	key      jwt.Keyfunc
}

func NewVerifier(
	ctx context.Context,
	issuer string,
	audience string,
	jwksURL string,
) (*Verifier, error) {
	if issuer == "" || audience == "" || jwksURL == "" {
		return nil, errors.New("issuer, audience and JWKS URL are required")
	}

	keys, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, err
	}

	return &Verifier{
		issuer:   issuer,
		audience: audience,
		key:      keys.Keyfunc,
	}, nil
}

func (v *Verifier) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		v.key,
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.Subject == "" ||
		(claims.AuthorizedParty != "idenova-dashboard" &&
			claims.AuthorizedParty != "idenova-services") {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

type contextKey struct{}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(*Claims)
	return claims, ok
}

func (v *Verifier) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.Fields(header)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := v.Verify(parts[1])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if len(roles) > 0 {
				allowed := false
				for _, role := range roles {
					if claims.AuthorizedParty == "idenova-dashboard" &&
						claims.HasRole(role) {
						allowed = true
						break
					}
				}
				if !allowed {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			ctx := context.WithValue(r.Context(), contextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireService accepts only tokens issued to the internal services client.
// Machine tokens cannot use this middleware to acquire human permissions.
func (v *Verifier) RequireService(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := v.Verify(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.AuthorizedParty != "idenova-services" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRolesOrService permits human users with the specified roles,
// or authenticated internal services, on explicitly selected endpoints.
func (v *Verifier) RequireRolesOrService(
	roles ...string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := v.Verify(parts[1])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			allowed := claims.AuthorizedParty == "idenova-services"
			if claims.AuthorizedParty == "idenova-dashboard" {
				for _, role := range roles {
					if claims.HasRole(role) {
						allowed = true
						break
					}
				}
			}

			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
