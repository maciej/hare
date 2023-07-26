package main

import (
	"embed"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

//go:embed static/*
var staticFS embed.FS

// See https://github.com/golang/go/issues/44854 for why we use a hardcoded modTime
var staticModTime = time.Now()

func serveStatic(w http.ResponseWriter, r *http.Request, file string) {
	f, err := staticFS.Open(filepath.Clean(filepath.Join("static", file)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, file, staticModTime, f.(io.ReadSeeker))
}
