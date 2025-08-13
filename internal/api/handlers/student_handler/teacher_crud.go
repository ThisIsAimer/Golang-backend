package students

import (
	"fmt"
	"net/http"
)

func GetStudentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Get")
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Get")
}

//post ----------------------------------------------------------------------------------
func PostStudentsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Post")
}

//put -----------------------------------------------------------------------------------
func PutStudentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Put")
}

//patch ------------------------------------------------------------------------------
func PatchStudentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Patch")
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Patch")
}

//delete -----------------------------------------------------------------------------
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Delete")
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "accessed : Students. with: Delete")
}