package teachers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"simpleapi/internal/models"
	teacherdb "simpleapi/internal/repositories/sql/teachersdb"
)

func PostTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if fieldIsEmpty(newTeachers) {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	newTeachers, err = teacherdb.PostTeachersDBHandler(w, newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(newTeachers),
		Data:   newTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func fieldIsEmpty(models []models.Teacher) bool {
	for _, value := range models {
		element := reflect.ValueOf(value)
		for i := range element.NumField() {
			if element.Field(i).Kind() == reflect.String && element.Field(i).String() == "" {
				return true
			}
		}
	}

	return false
}
