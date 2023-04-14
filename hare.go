package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "hare",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "addr", Usage: "Address to listen on", Value: "localhost:3200"},
		},
		Action: start,
	}

	_ = app.Run(os.Args)
}

func start(cCtx *cli.Context) error {
	mux := chi.NewMux()

	mux.Get("/headers", headersHandler)

	fmt.Printf("HARE starting. Listening on %s\n", cCtx.String("addr"))

	err := http.ListenAndServe(cCtx.String("addr"), mux)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error on shutdown: %s\n", err)
	}

	return nil
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r.Header)
}
