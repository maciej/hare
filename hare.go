package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
)

var addr string
var maxBodySize int64
var routeEnabled bool

func main() {
	app := &cli.App{
		Name: "hare",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "addr", Usage: "Address to listen on", Value: "localhost:3200", EnvVars: []string{"HARE_ADDR"}, Destination: &addr},
			&cli.Int64Flag{Name: "max-body-size", Usage: "Max body size", Value: 16 * 1024 * 1024, EnvVars: []string{"HARE_MAX_BODY_SIZE"}, Destination: &maxBodySize,
				Action: func(context *cli.Context, v int64) error {
					if v <= 0 {
						return errors.New("non-positive max-body-size")
					}
					return nil
				},
			},
			&cli.BoolFlag{Name: "route-enabled", Usage: "route-enabled", Value: false, EnvVars: []string{"HARE_ROUTE_ENABLED"}, Destination: &routeEnabled},
		},
		Action: start,
	}

	_ = app.Run(os.Args)
}

func start(cCtx *cli.Context) error {
	mux := chi.NewMux()

	mux.Get("/headers", headersHandler)
	mux.Get("/set-cookie", setCookieHandler)
	mux.Get("/hello", helloHandler)
	mux.Post("/body", bodyHandler)

	if routeEnabled {
		mux.Get("/route", routeHandler)
	}

	fmt.Printf("HARE starting. Listening on %s\n", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error on shutdown: %s\n", err)
	}

	return nil
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	route := r.Header.Get("HARE-Route")
	if route == "" {
		http.Error(w, "route missing", http.StatusBadRequest)
		return
	}

	targetReq, _ := http.NewRequest(http.MethodGet, route, nil)
	targetReq = targetReq.WithContext(r.Context())

	for k, vs := range r.Header {
		if len(vs) == 0 {
			continue
		}
		r.Header.Set(k, vs[0])

		for _, v := range vs[1:] {
			r.Header.Add(k, v)
		}
	}

	resp, err := http.Get(route)
	if err != nil {
		http.Error(w, "routing error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	copyHeader(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)

	_, _ = io.Copy(w, resp.Body)
}

func copyHeader(src, dst http.Header) {
	for k, vs := range src {
		if len(vs) == 0 {
			continue
		}
		dst.Set(k, vs[0])

		for _, v := range vs[1:] {
			src.Add(k, v)
		}
	}
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	var reqBuf bytes.Buffer
	var respContentLength int

	if r.ContentLength > 0 && r.ContentLength < maxBodySize {
		respContentLength = int(r.ContentLength)
	} else if r.ContentLength >= maxBodySize {
		respContentLength = int(maxBodySize)
	}

	if respContentLength > 0 {
		reqBuf.Grow(respContentLength)
		w.Header().Set("content-length", strconv.Itoa(respContentLength))
	}

	_, err := reqBuf.ReadFrom(http.MaxBytesReader(w, r.Body, maxBodySize))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := reqBuf.WriteTo(w)

	log.Printf("/body: %d bytes written, request content-length: %d, err: %v", n, r.ContentLength, err)
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
