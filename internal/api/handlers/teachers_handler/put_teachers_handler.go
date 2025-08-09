package teachers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/teachersdb"
)

func PutTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}

	updatedTeacher, existingTeacher, err := teacherdb.PutTeacherDBHandler(w, id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Status       string `json:"status"`
		OldEntry     models.Teacher
		UpdatedEntry models.Teacher
	}{
		Status:       "success",
		OldEntry:     existingTeacher,
		UpdatedEntry: updatedTeacher,
	}

	json.NewEncoder(w).Encode(responce)

}
