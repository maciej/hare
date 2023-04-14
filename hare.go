package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewMux()

	mux.Get("/headers", headersHandler)

	err := http.ListenAndServe("localhost:3200", mux)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error on shutdown: %s\n", err)
	}
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r.Header)
}
