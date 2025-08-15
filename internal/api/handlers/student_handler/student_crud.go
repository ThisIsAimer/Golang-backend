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

// get------------------------------------------------------------------------------------------------------
func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid id")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	student, err := studentdb.GetStudentDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return

	}
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	validTags := getModelTags(models.Student{})

	studentList, err := studentdb.GetStudentsDBHandler(r, validTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// post ----------------------------------------------------------------------------------------
func PostStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&students)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	modleTags := getModelTags(models.Student{})

	students, err = studentdb.PostStudentsDBHandler(modleTags, students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(students),
		Data:   students,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		myErr := utils.ErrorHandler(err, "Failed to encode response")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

// put -----------------------------------------------------------------------------------
func PutStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid id")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}
	var student models.Student

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&student)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	entries := getModelTags(models.Student{})

	student, existingStudent, err := studentdb.PutStudentDBHandler(id, student, entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status     string         `json:"status"`
		OldTeacher models.Student `json:"old_teacher"`
		NewTeacher models.Student `json:"new_teacher"`
	}{
		Status:     "Success",
		OldTeacher: existingStudent,
		NewTeacher: student,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		myErr := utils.ErrorHandler(err, "Failed to encode response")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

// patch ------------------------------------------------------------------------------
func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var updates map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}

	err = studentdb.PatchStudentDBHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}
	err = studentdb.PatchStudentsDBHandler(updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// delete -----------------------------------------------------------------------------
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	err = studentdb.DeleteStudentDBHandler(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responce := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "teacher successfully deleted",
		Id:     id,
	}

	err = json.NewEncoder(w).Encode(responce)
	
	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "accessed : Students. with: Delete")
}
