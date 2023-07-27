package main

import (
	"embed"
	"hare/middleware"
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

	rs := f.(io.ReadSeeker)

	etag, err := middleware.GenEtagFromReader(rs)
	if err == nil { // set ETag if possible
		w.Header().Set("Etag", etag)
	}

	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, file, staticModTime, rs)
}
