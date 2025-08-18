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

	// teachers-----------------------------------------------------------------------------------
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

	// students------------------------------------------------------------------------------------
	mux.HandleFunc("GET /students", students.GetStudentsHandler)
	mux.HandleFunc("POST /students", students.PostStudentsHandler)
	mux.HandleFunc("PATCH /students", students.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", students.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", students.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", students.PutStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", students.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", students.DeleteStudentHandler)
	

	// execs --------------------------------------------------------------------------------------------
	mux.HandleFunc("GET /execs", execs.ExecsRoute)
	mux.HandleFunc("POST /execs", execs.ExecsRoute)
	mux.HandleFunc("PATCH /execs", execs.ExecsRoute)
	mux.HandleFunc("DELETE /execs", execs.ExecsRoute)

	mux.HandleFunc("GET /execs/{id}", execs.ExecsRoute)
	mux.HandleFunc("PATCH /execs/{id}", execs.ExecsRoute)
	mux.HandleFunc("DELETE /execs/{id}", execs.ExecsRoute)

	mux.HandleFunc("POST /execs/login", execs.ExecsRoute)
	mux.HandleFunc("POST /execs/logout", execs.ExecsRoute)
	mux.HandleFunc("POST /execs/login/forgotpassword", execs.ExecsRoute)
	mux.HandleFunc("POST/execs/{id}/updatepassword", execs.ExecsRoute)
	mux.HandleFunc("POST /execs/login/resetpassword/reset/{resetcode}", execs.ExecsRoute)

	
	return mux
}