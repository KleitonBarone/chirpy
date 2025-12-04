package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/KleitonBarone/chirpy/internal/config"
	"github.com/KleitonBarone/chirpy/internal/database"
	"github.com/KleitonBarone/chirpy/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	const port = "8080"
	const staticDir = "."

	mux := http.NewServeMux()

	apiCfg := config.NewApiConfig(dbQueries, platform, jwtSecret)

	// Health check
	mux.HandleFunc("GET /api/healthz", handler.HandlerHealth)

	// Static files
	staticFiles := http.FileServer(http.Dir(staticDir))
	staticHandler := http.StripPrefix("/app", staticFiles)
	mux.Handle("/app/", handler.MiddlewareMetricsInc(apiCfg, staticHandler))

	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", handler.FileServerHitsHandler(apiCfg))
	mux.HandleFunc("POST /admin/reset", handler.FileServerHitsResetHandler(apiCfg))

	// Chirp endpoints
	mux.HandleFunc("POST /api/chirps", handler.CreateChirpHandler(apiCfg))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handler.GetChirpHandler(apiCfg))
	mux.HandleFunc("GET /api/chirps", handler.GetChirpsHandler(apiCfg))

	// User endpoints
	mux.HandleFunc("POST /api/users", handler.CreateUserHandler(apiCfg))
	mux.HandleFunc("POST /api/login", handler.LoginHandler(apiCfg))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
