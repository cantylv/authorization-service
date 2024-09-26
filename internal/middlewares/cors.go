// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"net/http"
)

// CORS (Cross-Origin Resource Sharing). Enabling communication with different services.
func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Need for postman | in real life for product version we should establish domain names instead of "*".
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS, HEAD")
		// Preflight-request processing.
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
	})
}
