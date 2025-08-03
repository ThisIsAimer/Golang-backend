package home

import (
	"fmt"
	"net/http"
)

// http methods are get, post, put, patch, delete

func HomeRoute(w http.ResponseWriter, r *http.Request) {

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
