package middlewares

import (
	"fmt"
	"net/http"
)

// cross-origine resource sharing
func Cors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println("Origin:", origin)

		if origin == "https://localhost:3000" {
			fmt.Println("access allowed")
		} else {
			http.Error(w, "not allowed by Cors", http.StatusForbidden)
			return
		}

		// Set other CORS headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		next.ServeHTTP(w, r)
	})

}
