package middleware

import (
	"net/http"
)

// CORS handles Cross-Origin Resource Sharing logic to ensure the proxy can be accessed
// by allowed origins (like your Next.js portfolio) and handles OPTIONS preflight requests.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// If allowedOrigin is "*", we allow any origin. Otherwise, we only set the header
			// if the origin matches our allowed list. For simplicity and security, we enforce it.
			if allowedOrigin == "*" || origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// Intercept preflight OPTIONS request and return immediately.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
