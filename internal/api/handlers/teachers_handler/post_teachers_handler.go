package teachers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"simpleapi/internal/models"
	teacherdb "simpleapi/internal/repositories/sql/teachersdb"
	"simpleapi/pkg/utils"
)

func PostTeachersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	// used to discard unknown fields
	decoder.DisallowUnknownFields()

	var newTeachers []models.Teacher
	err := decoder.Decode(&newTeachers)
	if err != nil {
		myError :=  utils.ErrorHandler(err,"invalid request body")
		http.Error(w, myError.Error(), http.StatusBadRequest)
		return
	}

	for _, teacher := range newTeachers {
		err = fieldIsEmpty(teacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

func fieldIsEmpty(model any) error {
	element := reflect.ValueOf(model)
	for i := range element.NumField() {
		if element.Field(i).Kind() == reflect.String && element.Field(i).String() == "" {
			return utils.ErrorHandler(errors.New("user has not provided all fields"), "all fields required")
		}
	}

	return nil
}
