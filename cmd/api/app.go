package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	mid "simpleapi/internal/api/middlewares"
)

// for now
type Teacher struct {
	ID         int
	FirestName string
	LastName   string
	Class      string
	Subject    string
}

var (
	teachers = make(map[int]Teacher)
	mutex = &sync.Mutex{}
	nextId = 1
)

func init(){
	teachers[nextId] = Teacher{
		ID: nextId,
		FirestName: "Rudra",
		LastName: "ABC",
		Class: "6A",
		Subject: "math",
	}
	nextId++

	teachers[nextId] = Teacher{
		ID: nextId,
		FirestName: "Rudrina",
		LastName: "ABC",
		Class: "10A",
		Subject: "computer",
	}

}


func getTeachersHandler(w http.ResponseWriter, r *http.Request){

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
		CheckQuery:              true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		WhiteList:               []string{"allowedParam", "sortOrder", "sortBy", "name", "age", "class"},
	}

	hppMiddleware := mid.Hpp(*hppSettings)

	// secureMux := mid.Cors(rateLimiter.Middleware(mid.ResponseTime(mid.SecurityHeaders(mid.CompMiddleware(hppMiddleware(mux))))))
	// secureMux := applyMiddlewares(mux,hppMiddleware,mid.CompMiddleware,mid.SecurityHeaders,mid.ResponseTime,rateLimiter.Middleware,mid.Cors)
	secureMux := applyMiddlewares(mux, hppMiddleware, rateLimiter.Middleware) // for now faster processing

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

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {

	for _, v := range middlewares {
		handler = v(handler)
	}
	return handler
}
