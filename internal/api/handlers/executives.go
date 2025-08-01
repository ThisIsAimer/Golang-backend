package handlers

import (
	"fmt"
	"net/http"
)

func ExecsRoute(w http.ResponseWriter, r *http.Request) {
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