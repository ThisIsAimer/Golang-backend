package teachers

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"simpleapi/internal/models"
	teacherdb "simpleapi/internal/repositories/sql/teachersdb"
)

func PostTeachersHandler(w http.ResponseWriter, r *http.Request) {

	// new decoder only reads from reader once
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var verifyJson []map[string]any
	err = json.Unmarshal(bodyBytes, &verifyJson)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if keysInvalid(verifyJson) {
		http.Error(w, "invalid json keys included", http.StatusBadRequest)
		return
	}

	var newTeachers []models.Teacher
	err = json.Unmarshal(bodyBytes, &newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, teacher := range newTeachers{
		if fieldIsEmpty(teacher) {
		http.Error(w, "all fields are required", http.StatusBadRequest)
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

func fieldIsEmpty(model any) bool {
	element := reflect.ValueOf(model)
	for i := range element.NumField() {
		if element.Field(i).Kind() == reflect.String && element.Field(i).String() == "" {
			return true
		}
	}
	

	return false
}

func keysInvalid(data []map[string]any) bool {
	for _, teacher := range data {
		for k := range teacher {
			if checkValidKey(k) {
				continue
			} else {
				return true
			}
		}
	}

	return false
}

func checkValidKey(key string) bool {
	modelType := reflect.TypeOf(models.Teacher{})
	validKey := make(map[string]bool)

	for i := range modelType.NumField() {
		modelTag := strings.TrimSuffix(modelType.Field(i).Tag.Get("json"), ",omitempty")

		validKey[modelTag] = true
	}
	return validKey[key]
}
