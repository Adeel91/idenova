package main

import (
	"context"
	"log"
	"net/http"
	"os"

	httpadapter "github.com/Adeel91/idenova/services/verification/internal/adapters/http"
	"github.com/Adeel91/idenova/services/verification/internal/adapters/postgres"
	"github.com/Adeel91/idenova/services/verification/internal/adapters/postgres/db"
	"github.com/Adeel91/idenova/services/verification/internal/application"
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
	handler := httpadapter.NewHandler(createUseCase, getUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /verifications", handler.CreateVerification)
	mux.HandleFunc("GET /verifications/{id}", handler.GetVerification)

	log.Println("Verification Service running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("HTTP server stopped: %v", err)
		return
	}
}
