package teachers

import (
	"fmt"
	"net/http"
)

func TeachersRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Teachers route")

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandler(w, r)
	case http.MethodPut:
		PutTeachersHandler(w, r)
	case http.MethodPatch:
		PatchTeachersHandler(w, r)
	case http.MethodDelete:
		DeleteTeachersHandler(w, r)
	default:
		fmt.Fprintln(w, "accessed : Teachers")

	}
}
