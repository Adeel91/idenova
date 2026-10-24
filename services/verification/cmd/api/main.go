package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"idenova/auth"
	"idenova/verification/internal/adapters/documenthttp"
	httpadapter "idenova/verification/internal/adapters/http"
	"idenova/verification/internal/adapters/postgres"
	"idenova/verification/internal/adapters/postgres/db"
	"idenova/verification/internal/application"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("database connection failed: %v", err)
		return
	}

	repo := postgres.NewRepository(db.New(pool))
	createUseCase := application.NewCreateVerification(repo)
	getUseCase := application.NewGetVerification(repo)
	listUseCase := application.NewListVerifications(repo)
	documentServiceURL := os.Getenv("DOCUMENT_SERVICE_URL")
	if documentServiceURL == "" {
		documentServiceURL = "http://127.0.0.1:8083"
	}

	serviceClient, err := auth.NewServiceHTTPClient()
	if err != nil {
		log.Printf("Failed to initialize service HTTP client: %v", err)
		return
	}
	documentReader := documenthttp.NewReader(documentServiceURL, serviceClient)

	updateUseCase := application.NewUpdateVerification(repo, documentReader)
	handler := httpadapter.NewHandler(createUseCase, getUseCase, listUseCase, updateUseCase)

	handler.SetIdempotentCreation(
		application.NewCreateVerificationIdempotent(
			postgres.NewIdempotencyStore(pool),
		),
	)

	issuer := "http://localhost:8085/realms/idenova"

	verifier, err := auth.NewVerifier(
		context.Background(),
		issuer,
		"idenova-api",
		issuer+"/protocol/openid-connect/certs",
	)
	if err != nil {
		log.Printf("Keycloak JWT verifier initialization failed: %v", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle(
		"POST /verifications",
		verifier.RequireRoles("ADMIN", "REVIEWER")(
			http.HandlerFunc(handler.CreateVerification),
		),
	)
	mux.HandleFunc("GET /verifications/{id}", handler.GetVerification)
	mux.HandleFunc("GET /verifications", handler.ListVerifications)
	mux.Handle(
		"PATCH /verifications/{id}/status",
		verifier.RequireRoles("ADMIN", "REVIEWER")(
			http.HandlerFunc(handler.UpdateVerification),
		),
	)

	log.Println("Verification Service running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("HTTP server stopped: %v", err)
		return
	}
}
