package middleware

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
)

func BodyEtag(next http.Handler) http.Handler {
	return bodyEtag(false)(next)
}

func WeakBodyEtag(next http.Handler) http.Handler {
	return bodyEtag(true)(next)
}

func bodyEtag(weak bool) func(next http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			brw := &bufferedResponseWriter{
				Buffer:         bytes.NewBuffer(nil),
				ResponseWriter: w,
			}
			next.ServeHTTP(brw, r)

			etag := GenEtag(brw.Bytes())

			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			var prefix = ""
			if weak {
				prefix = "W/"
			}

			w.Header().Set("Etag", prefix+GenEtag(brw.Bytes()))

			_, _ = io.Copy(brw.ResponseWriter, brw.Buffer) // TODO figure out error handling
		}

		return http.HandlerFunc(fn)
	}

	return mw
}

var _ http.ResponseWriter = (*bufferedResponseWriter)(nil)

type bufferedResponseWriter struct {
	*bytes.Buffer
	http.ResponseWriter
}

func (brw *bufferedResponseWriter) Write(buf []byte) (int, error) {
	return brw.Buffer.Write(buf)
}

func GenEtag(buf []byte) string {
	sum := sha1.Sum(buf)
	return encEtagVal(sum[:])
}

func GenEtagFromReader(reader io.Reader) (string, error) {
	sum := sha1.New()

	if _, err := io.Copy(sum, reader); err != nil {
		return "", err
	}

	return encEtagVal(sum.Sum(nil)), nil
}

func encEtagVal(buf []byte) string {
	return "\"" + base64.StdEncoding.EncodeToString(buf) + "\""
}
