package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	mid "simpleapi/internal/api/middlewares"
)

type user struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Place string `json:"place"`
}

// http methods are get, post, put, patch, delete

func homeRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: home")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "accessed : Home. with: Get")
	case http.MethodPost:
		fmt.Fprintln(w, "accessed : Home. with: Post")
		fmt.Println("quary:", r.URL.Query())
	
		r.ParseForm()

		// r.form includes quary parameters
		// to get just post form we use r.PostForm
		fmt.Println("form:", r.PostForm)
		fmt.Println("form:", r.PostForm.Get("allowedParam"))
	case http.MethodPut:
		fmt.Fprintln(w, "accessed : Home. with: Put")
	case http.MethodPatch:
		fmt.Fprintln(w, "accessed : Home. with: Patch")
	case http.MethodDelete:
		fmt.Fprintln(w, "accessed : Home. with: Delete")
	default:
		fmt.Fprintln(w, "accessed : Home")

	}

}

func teachersRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Teachers route")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "accessed : Teachers. with: Get")
	case http.MethodPost:
		fmt.Fprintln(w, "accessed : Teachers. with: Post")
	case http.MethodPut:
		fmt.Fprintln(w, "accessed : Teachers. with: Put")
	case http.MethodPatch:
		fmt.Fprintln(w, "accessed : Teachers. with: Patch")
	case http.MethodDelete:
		fmt.Fprintln(w, "accessed : Teachers. with: Delete")
	default:
		fmt.Fprintln(w, "accessed : Teachers")

	}
}

func studentsRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Students route")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "accessed : Students. with: Get")
	case http.MethodPost:
		fmt.Fprintln(w, "accessed : Students. with: Post")
	case http.MethodPut:
		fmt.Fprintln(w, "accessed : Students. with: Put")
	case http.MethodPatch:
		fmt.Fprintln(w, "accessed : Students. with: Patch")
	case http.MethodDelete:
		fmt.Fprintln(w, "accessed : Students. with: Delete")
	default:
		fmt.Fprintln(w, "accessed : Students")

	}
}

func execsRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Execs route")
	fmt.Println("method:", r.Method)

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "accessed : Executives. with: Get")
	case http.MethodPost:
		fmt.Fprintln(w, "accessed : Executives. with: Post")
	case http.MethodPut:
		fmt.Fprintln(w, "accessed : Executives. with: Put")
	case http.MethodPatch:
		fmt.Fprintln(w, "accessed : Executives. with: Patch")
	case http.MethodDelete:
		fmt.Fprintln(w, "accessed : Executives. with: Delete")
	default:
		fmt.Fprintln(w, "accessed : Executives")

	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeRoute)
	mux.HandleFunc("/teachers", teachersRoute)
	mux.HandleFunc("/students", studentsRoute)
	mux.HandleFunc("/execs", execsRoute)

	port := 3000

	key := `certificate\key.pem`
	cert := `certificate\certificate.pem`

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rateLimiter := mid.NewRateLimiter(5, time.Second*5)

	hppSettings := &mid.HppOptions{
		CheckQuery: true,
		CheckBody: true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		WhiteList: []string{"allowedParam"},
	}

	hppMiddleware := mid.Hpp(*hppSettings)

	secureMux := hppMiddleware(rateLimiter.Middleware(mid.CompMiddleware(mid.ResponseTime(mid.SecurityHeaders(mid.Cors(mux))))))

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server is running on port:", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}
