package router

import (
	"net/http"

	home "simpleapi/internal/api/handlers/home_handler"
	students "simpleapi/internal/api/handlers/student_handler"
	teachers "simpleapi/internal/api/handlers/teachers_handler"
	execs "simpleapi/internal/api/handlers/execs_handler"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home.HomeRoute)
	mux.HandleFunc("/teachers/", teachers.TeachersRoute)
	mux.HandleFunc("/students/", students.StudentsRoute)
	mux.HandleFunc("/execs/", execs.ExecsRoute)
	return mux
}
