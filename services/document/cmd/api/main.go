package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "idenova/document/internal/adapters/http"
	"idenova/document/internal/adapters/postgres"
	"idenova/document/internal/adapters/postgres/db"
	"idenova/document/internal/adapters/verificationhttp"
	"idenova/document/internal/application"

	"idenova/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	repo := postgres.NewRepository(db.New(pool))

	handler := httpadapter.NewHandler(
		application.NewCreateDocument(repo),
		application.NewGetDocument(repo),
		application.NewListDocuments(repo),
		application.NewUpdateDocument(repo),
	)

	serviceClient, err := auth.NewServiceHTTPClient()
	if err != nil {
		return err
	}
	handler.SetIdempotentCreation(
		application.NewCreateDocumentIdempotent(
			postgres.NewIdempotencyStore(
				pool,
				verificationhttp.NewReader("http://127.0.0.1:8080", serviceClient),
			),
		),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /documents", handler.CreateDocument)
	mux.HandleFunc("GET /documents/{id}", handler.GetDocument)
	mux.HandleFunc("GET /verifications/{id}/documents", handler.ListDocuments)
	mux.HandleFunc("PATCH /documents/{id}/status", handler.UpdateDocument)

	server := &http.Server{
		Addr:              "127.0.0.1:8083",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Document Service running on http://localhost:8083")

	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
