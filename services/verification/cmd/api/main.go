package main

import (
	"context"
	"log"
	"net/http"
	"os"

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
	updateUseCase := application.NewUpdateVerification(repo)
	handler := httpadapter.NewHandler(createUseCase, getUseCase, updateUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /verifications", handler.CreateVerification)
	mux.HandleFunc("GET /verifications/{id}", handler.GetVerification)
	mux.HandleFunc("PATCH /verifications/{id}/status", handler.UpdateVerification)

	log.Println("Verification Service running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("HTTP server stopped: %v", err)
		return
	}
}
