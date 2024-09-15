package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("no port configured")
	}
	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Fatalf("Server starting on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server")
	}
}
