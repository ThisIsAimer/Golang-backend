package router

import (
	"net/http"
	students "simpleapi/internal/api/handlers/student_handler"
)

func studentsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// students
	mux.HandleFunc("GET /students", students.GetStudentsHandler)
	mux.HandleFunc("POST /students", students.PostStudentsHandler)
	mux.HandleFunc("PATCH /students", students.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", students.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", students.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", students.PutStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", students.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", students.DeleteStudentHandler)

	return mux
}