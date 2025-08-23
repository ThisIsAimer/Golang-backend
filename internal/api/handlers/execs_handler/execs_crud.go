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
		Status string              `json:"status"`
		Count  int                 `json:"count"`
		Data   []models.BasicExecs `json:"data"`
	}{
		Status: "success",
		Count:  len(execsList),
		Data:   execsList,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// Post------------------------------------------------------------------------------------------------
func PostExecsHandler(w http.ResponseWriter, r *http.Request) {
	var execs []models.Execs
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&execs)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	execTags := getModelTags(models.Execs{})

	execs, err = execsdb.PostExecsDBHandler(execTags, execs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []models.Execs `json:"data"`
	}{
		Status: "Success",
		Count:  len(execs),
		Data:   execs,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		myErr := utils.ErrorHandler(err, "Failed to encode response")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

// Patch----------------------------------------------------------------------------------------------
func PatchExecHandler(w http.ResponseWriter, r *http.Request) {
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

	err = execsdb.PatchExecDBHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

}

// Delete----------------------------------------------------------------------------------------------
func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteExecsHandler(w http.ResponseWriter, r *http.Request) {

}

// Login----------------------------------------------------------------------------------------------
func LoginExecHandler(w http.ResponseWriter, r *http.Request) {

}

func LogoutExecHandler(w http.ResponseWriter, r *http.Request) {

}

// Passwords----------------------------------------------------------------------------------------------
func ForgetPassExecHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePassExecHandler(w http.ResponseWriter, r *http.Request) {

}

func ResetPassExecHandler(w http.ResponseWriter, r *http.Request) {

}
