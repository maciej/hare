package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "hare",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "addr", Usage: "Address to listen on", Value: "localhost:3200", EnvVars: []string{"HARE_ADDR"}},
		},
		Action: start,
	}

	_ = app.Run(os.Args)
}

func start(cCtx *cli.Context) error {
	addr := cCtx.String("addr")

	mux := chi.NewMux()

	mux.Get("/headers", headersHandler)
	mux.Get("/set-cookie", setCookieHandler)
	mux.Get("/hello", helloHandler)
	mux.Post("/body", bodyHandler)

	fmt.Printf("HARE starting. Listening on %s\n", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error on shutdown: %s\n", err)
	}

	return nil
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.Copy(w, r.Body)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello")
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "HARE-Hello", Value: "1", MaxAge: 3600})
	w.Header().Set("content-type", "text/plain")

	_, _ = fmt.Fprintln(w, "OK")
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	var renderJSON bool

	switch r.Header.Get("accept") {
	case "text/plain":
		renderJSON = false
	case "application/json", "text/json":
		renderJSON = true
	default:
		renderJSON = false
	}

	if !renderJSON {
		w.Header().Set("content-type", "text/plain")
		_ = r.Header.Write(w)
	} else {
		w.Header().Set("content-type", "application/json")

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(r.Header)
	}

}
