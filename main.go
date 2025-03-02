package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const staticDir = "."

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlerHealth)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	staticFiles := http.FileServer(http.Dir(staticDir))
	staticHandler := http.StripPrefix("/app", staticFiles)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(staticHandler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.fileServerHitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.fileServerHitsResetHandler)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
