package main

import (
	"log"
	"net/http"

	httpadapter "github.com/Adeel91/idenova/services/verification/internal/adapters/http"
	"github.com/Adeel91/idenova/services/verification/internal/adapters/memory"
	"github.com/Adeel91/idenova/services/verification/internal/application"
)

func main() {
	repo := memory.NewRepository()
	createUseCase := application.NewCreateVerification(repo)
	getUseCase := application.NewGetVerification(repo)
	handler := httpadapter.NewHandler(createUseCase, getUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /verifications", handler.CreateVerification)
	mux.HandleFunc("GET /verifications/{id}", handler.GetVerification)

	log.Println("Verification Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
