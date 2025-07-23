package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

type user struct{
	Name string `json:"name"`
	Age int `json:"age"`
	Place string `json:"place"`
}


// http methods are get, post, put, patch, delete

func homeRoute(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: home")
	fmt.Println("method:", r.Method)

	switch r.Method{
	case http.MethodGet:
		// for a specific ID, it will be routed to /teachers/90 or smth
		urlPath := strings.TrimPrefix(r.URL.Path, "/")
		userId := strings.TrimSuffix(urlPath,"/")

		if userId != ""{
			fmt.Println("id is:",userId)
		}

		//------------handling quary-----------------------------------
		if len(r.URL.Query()) > 0{

			quaryParams := r.URL.Query()

			for key := range quaryParams{
				fmt.Printf("key %v : %v \n", key, quaryParams.Get(key))
			}
		}


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

func teachersRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Teachers route")
	fmt.Println("method:", r.Method)

	switch r.Method{
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

func studentsRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Students route")
	fmt.Println("method:", r.Method)

	switch r.Method{
	case http.MethodGet:
		// for a specific ID, it will be routed to /teachers/90 or smth
		urlPath := strings.TrimPrefix(r.URL.Path, "/students/")
		userId := strings.TrimSuffix(urlPath,"/")

		if userId != ""{
			fmt.Println("id is:",userId)
		}

		//------------handling quary-----------------------------------
		if len(r.URL.Query()) > 0{

			quaryParams := r.URL.Query()

			for key := range quaryParams{
				fmt.Printf("key %v : %v \n", key, quaryParams.Get(key))
			}
		}


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

func execsRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Execs route")
	fmt.Println("method:", r.Method)

	switch r.Method{
	case http.MethodGet:
		// for a specific ID, it will be routed to /teachers/90 or smth
		urlPath := strings.TrimPrefix(r.URL.Path, "/execs/")
		userId := strings.TrimSuffix(urlPath,"/")

		if userId != ""{
			fmt.Println("id is:",userId)
		}

		//------------handling quary-----------------------------------
		if len(r.URL.Query()) > 0{

			quaryParams := r.URL.Query()

			for key := range quaryParams{
				fmt.Printf("key %v : %v \n", key, quaryParams.Get(key))
			}
		}


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

func main(){

	http.HandleFunc("/", homeRoute)
	http.HandleFunc("/teachers/", teachersRoute)
	http.HandleFunc("/students/", studentsRoute)
	http.HandleFunc("/execs/", execsRoute)


	port := 3000

	key := `certificate\key.pem`
	cert := `certificate\certificate.pem`

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}


	server := &http.Server{
		Addr: fmt.Sprintf(":%d",port),
		Handler: nil,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server is running on port:", port)


	err := server.ListenAndServeTLS(cert,key)
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}