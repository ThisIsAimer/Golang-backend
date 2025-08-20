package execs

import (
	"encoding/json"
	"net/http"
	
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/execsdb"
	"simpleapi/pkg/utils"
	"strconv"
)

func GetExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid id")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	student, err := execsdb.GetExecDBHandler(id)
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

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	validTags := getModelTags(models.Student{})

	execsList, err := execsdb.GetExecsDBHandler(r, validTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []models.Execs `json:"data"`
	}{
		Status: "success",
		Count:  len(execsList),
		Data:   execsList,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
