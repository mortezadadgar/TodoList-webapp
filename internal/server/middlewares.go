package server

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func (s server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.u.IsAuthenticated(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s server) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				s.l.Error("%s\nSTACK: %s", err, debug.Stack())
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// GET http://localhost:4050/todo/id HTTP/1.1" from 127.0.0.1:37854 - 000 0B in 363.701µs
func (s server) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		defer func() {
			s.l.Info("%s %s %s from %s in %dms",
				r.Method,
				r.URL.Path,
				r.Proto,
				r.Host,
				time.Since(now))
		}()

		next.ServeHTTP(w, r)
	})
}

func (s server) restrictStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
