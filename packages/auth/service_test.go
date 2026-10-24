package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMachineAuthorization(t *testing.T) {
	verifier, key := testVerifier(t)

	token := func(client, role string) string {
		return signTestToken(
			t, key, role,
			verifier.issuer,
			verifier.audience,
			client,
			time.Now().Add(time.Hour),
		)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	tests := []struct {
		name    string
		handler http.Handler
		token   string
		want    int
	}{
		{
			"service allowed on service endpoint",
			verifier.RequireService(next),
			token("idenova-services", ""),
			204,
		},
		{
			"dashboard rejected on service endpoint",
			verifier.RequireService(next),
			token("idenova-dashboard", "ADMIN"),
			403,
		},
		{
			"service cannot impersonate ADMIN",
			verifier.RequireRoles("ADMIN")(next),
			token("idenova-services", "ADMIN"),
			403,
		},
		{
			"service allowed on shared read endpoint",
			verifier.RequireRolesOrService("ADMIN", "REVIEWER")(next),
			token("idenova-services", ""),
			204,
		},
		{
			"reviewer allowed on shared read endpoint",
			verifier.RequireRolesOrService("ADMIN", "REVIEWER")(next),
			token("idenova-dashboard", "REVIEWER"),
			204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)

			rec := httptest.NewRecorder()
			tt.handler.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, rec.Code)
			}
		})
	}
}
