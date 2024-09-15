package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MishaYanov/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("no port configured")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("no DB configured")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Connected to database successfully!")

	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/err", handleError)

	v1Router.Post("/users", apiCfg.HandleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.HandleGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.HandleCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandleGetFeeds)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server")
	}
}
