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

	// teachers
	mux.HandleFunc("GET /teachers/", teachers.GetTeacherHandler)
	mux.HandleFunc("POST /teachers/", teachers.PostTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", teachers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", teachers.DeleteTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}", teachers.GetTeachersHandler)
	mux.HandleFunc("PUT /teachers/{id}", teachers.PutTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", teachers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", teachers.DeleteTeacherHandler)


	mux.HandleFunc("/students/", students.StudentsRoute)

	mux.HandleFunc("/execs/", execs.ExecsRoute)
	return mux
}
