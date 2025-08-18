package router

import (
	"net/http"

	execs "simpleapi/internal/api/handlers/execs_handler"
	home "simpleapi/internal/api/handlers/home_handler"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", home.HomeRoute)
	mux.HandleFunc("/execs", execs.ExecsRoute)

	tRouter :=  teachersRouter()
	tRouter.Handle("/", studentsRouter())
	mux.Handle("/", tRouter)

	return mux
}
