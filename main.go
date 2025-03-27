package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/KleitonBarone/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
}

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

	mux.HandleFunc("GET /api/healthz", handlerHealth)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	staticFiles := http.FileServer(http.Dir(staticDir))
	staticHandler := http.StripPrefix("/app", staticFiles)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(staticHandler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.fileServerHitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.fileServerHitsResetHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
