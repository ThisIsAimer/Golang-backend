package router

import (
	"net/http"
	teachers "simpleapi/internal/api/handlers/teachers_handler"
)

func teachersRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// teachers
	mux.HandleFunc("GET /teachers", teachers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", teachers.PostTeachersHandler)
	mux.HandleFunc("PATCH /teachers", teachers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", teachers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", teachers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", teachers.PutTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", teachers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", teachers.DeleteTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", teachers.GetStudentsByTeacherId)
	mux.HandleFunc("GET /teachers/{id}/studentcount", teachers.GetStudentCountByTeacherId)
	return mux
}