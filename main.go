package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const staticDir = "."

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(staticDir))

	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
