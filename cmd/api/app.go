package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		fmt.Fprintln(w, "accessed : Home. with: Get")
	case http.MethodPost:
		fmt.Fprintln(w, "accessed : Home. with: Post")

		// parsing form in post request (necessary for x-www-form-urlencoded)

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "erroring parsing form", http.StatusBadRequest)
			return
		}

		//extracting the values

		response := make(map[string]any)

		for key, value := range r.Form{
			response[key] = value[0]
		}

		if len(response) > 0{
			fmt.Println(response)
		}

		//----------------------------------------------------------------------------------
		// raw response

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w,"error reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		println("raw body", string(body))

		var myUser user

		err = json.Unmarshal(body, &myUser)
		if err != nil {
			fmt.Println("error unmarshalling json data", err)
			return
		}

		fmt.Println("myUser:", myUser)
		fmt.Println("recieved userName:", myUser.Name)


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
	http.HandleFunc("/teachers", teachersRoute)
	http.HandleFunc("/students", studentsRoute)
	http.HandleFunc("/execs", execsRoute)


	port := 3000

	server := &http.Server{
		Addr: fmt.Sprintf(":%d",port),
		Handler: nil,
	}

	fmt.Println("server is running on port:", port)


	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}