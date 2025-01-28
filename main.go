package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const staticDir = "."
	const assetsDir = "assets"

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(staticDir))
	fsAssets := http.FileServer(http.Dir(assetsDir))

	mux.Handle("/", fs)
	mux.Handle("/assets", fsAssets)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
