package execs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	var updates []map[string]any

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}

	err = execsdb.PatchExecsDBHandler(updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Delete----------------------------------------------------------------------------------------------
func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid exec id", http.StatusBadRequest)
		return
	}

	err = execsdb.DeleteExecDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responce := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "exec successfully deleted",
		Id:     id,
	}

	err = json.NewEncoder(w).Encode(responce)

	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

func DeleteExecsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []any
	var intIds []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}

	for _, id := range ids {
		switch v := id.(type) {
		case float64:
			intIds = append(intIds, int(v))
		case int:
			intIds = append(intIds, v)
		case string:
			convID, err := strconv.Atoi(v)

			if err != nil {
				myErr := utils.ErrorHandler(err, "invalid id")
				http.Error(w, myErr.Error(), http.StatusBadRequest)
				return
			}

			intIds = append(intIds, convID)
		default:
			myErr := utils.ErrorHandler(fmt.Errorf("default type activated"), "invalid id")
			http.Error(w, myErr.Error(), http.StatusBadRequest)
			return

		}
	}
	err = execsdb.DeleteExecsDBHandler(intIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responce := struct {
		Status string `json:"status"`
		Ids    []int  `json:"ids"`
	}{
		Status: "students successfully deleted",
		Ids:    intIds,
	}

	err = json.NewEncoder(w).Encode(responce)

	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

// Login----------------------------------------------------------------------------------------------
func LoginExecHandler(w http.ResponseWriter, r *http.Request) {

	// data validatiuon
	var req models.Execs

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&req)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.Password == "" {
		myErr := utils.ErrorHandler(err, "username and password are required")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	givenPass := req.Password

	// search if user exists
	req, err = execsdb.LoginExecDBHandler(req.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// is user active
	if req.InactiveStatus {
		myErr := utils.ErrorHandler(fmt.Errorf("user inactive"), "user is inactive")
		http.Error(w, myErr.Error(), http.StatusForbidden)
		return
	}

	// verify password
	parts := strings.Split(req.Password, ".")
	if len(parts) != 2 {
		myErr := utils.ErrorHandler(fmt.Errorf("invalid encode hash format"), "Password must be reset")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
	}

	saltBase64 := parts[0]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		myErr := utils.ErrorHandler(err, "error decoding salt")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	givenPass, err = passEncoder(givenPass, salt)

	if givenPass != req.Password {
		myErr := utils.ErrorHandler(fmt.Errorf("password doesnt match"), "incorrect password")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	// generate token
	tokenString := "abc"

	// send token as response or a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
	})
	
	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "test",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

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
