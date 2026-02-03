package middleware

import (
	"net/http"
	"os"
)

func NewCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origins := []string{
			"http://fifsky.com",
			"https://www.fifsky.com",
			"http://www.fifsky.com",
			"https://fifsky.com",
			"https://windiness.fifsky.com",
		}

		if os.Getenv("APP_ENV") == "local" {
			origins = []string{"*"}
		}
		origin := r.Header.Get("Origin")
		if len(origin) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		host := r.Host

		if origin == "http://"+host || origin == "https://"+host {
			next.ServeHTTP(w, r)
			return
		}

		allowedOrigin := getOrigin(origin, origins)
		if allowedOrigin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "43200")

		reqMethod := r.Header.Get("Access-Control-Request-Method")
		if reqMethod != "" {
			w.Header().Set("Access-Control-Allow-Methods", reqMethod)
		} else {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}

		reqHeaders := r.Header.Get("Access-Control-Request-Headers")
		if reqHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Length,Content-Type,Access-Token,X-Requested-With")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getOrigin(requestedOrigin string, allowedOrigins []string) string {
	for _, origin := range allowedOrigins {
		if origin == "*" {
			return requestedOrigin
		}
		if requestedOrigin == origin {
			return requestedOrigin
		}
	}
	return ""
}
