package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	mid "simpleapi/internal/api/middlewares"
)

// for now
type Teacher struct {
	ID         int
	FirstName string
	LastName   string
	Class      string
	Subject    string
}

var (
	teachers = make(map[int]Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = Teacher{
		ID:         nextId,
		FirstName: "Rudra",
		LastName:   "Sivdev",
		Class:      "6A",
		Subject:    "math",
	}
	nextId++

	teachers[nextId] = Teacher{
		ID:         nextId,
		FirstName: "Rudrina",
		LastName:   "ShivDev",
		Class:      "10B",
		Subject:    "computer",
	}

	nextId++

	teachers[nextId] = Teacher{
		ID:         nextId,
		FirstName: "Tanjiro",
		LastName:   "Kamado",
		Class:      "all",
		Subject:    "Dance",
	}

	nextId++

	teachers[nextId] = Teacher{
		ID:         nextId,
		FirstName: "Zenitsu",
		LastName:   "Agatsuma",
		Class:      "8C",
		Subject:    "Science",
	}

	nextId++

	teachers[nextId] = Teacher{
		ID:         nextId,
		FirstName: "Inosuke",
		LastName:   "Hashibira",
		Class:      "5D",
		Subject:    "Sports",
	}

}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	teacherList := make([]Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName){
			teacherList = append(teacherList, teacher)
		}
	}

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		getTeachersHandler(w, r)
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
		WhiteList:               []string{"allowedParam", "sortOrder", "sortBy", "name", "age", "class", "first_name", "last_name"},
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
