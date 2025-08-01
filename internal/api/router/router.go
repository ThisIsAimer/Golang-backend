package router

import (
	"net/http"
	"simpleapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeRoute)
	mux.HandleFunc("/teachers/", handlers.TeachersRoute)
	mux.HandleFunc("/students/", handlers.StudentsRoute)
	mux.HandleFunc("/execs/", handlers.ExecsRoute)
	return mux
}
