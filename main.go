package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthzHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, "OK")
}

func main() {
	const port = "8080"
	const staticDir = "."

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler)

	staticFiles := http.FileServer(http.Dir(staticDir))
	mux.Handle("/app/", http.StripPrefix("/app", staticFiles))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
