package students

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/studentdb"
	"simpleapi/pkg/utils"
)

func GetStudentHandler(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		myErr := utils.ErrorHandler(err,"invalid id")
		http.Error(w, myErr.Error(),http.StatusBadRequest)
		return
	}

	student, err :=  studentdb.GetStudentDBHandler(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return

	}
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request){
	validTags := getModelTags(models.Student{})

	studentList , err := studentdb.GetStudentsDBHandler(r,validTags)
	if err != nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
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