package execs

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	defer r.Body.Close()

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
	defer r.Body.Close()

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
	defer r.Body.Close()

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
	defer r.Body.Close()

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

	err = utils.VerifyPassword(givenPass, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate token
	tokenString, err := utils.SignToken(req.ID, req.UserName, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send token as response or a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "test",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	err = json.NewEncoder(w).Encode(responce)

	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

func LogoutExecHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("Message: logged out successfully"))
}

// Passwords----------------------------------------------------------------------------------------------
func ForgetPassExecHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "please enter an email", http.StatusBadRequest)
		return
	}

	err = execsdb.ForgotPasswordDBHandler(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Status string `json:"status"`
	}{
		Status: fmt.Sprintf("Sent reset link to email : %s", req.Email),
	}

	err = json.NewEncoder(w).Encode(responce)

	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

func UpdatePassExecHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid execs id")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	var req models.UpdatePasswordRequest

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.DisallowUnknownFields()

	err = decoder.Decode(&req)
	if err != nil {
		myErr := utils.ErrorHandler(err, "invalid json body")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		myErr := utils.ErrorHandler(fmt.Errorf("current or new passwords are unfilled"), "please enter password")
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	err = execsdb.UpdatePassExecDBHandler(userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Message string `json:"message"`
	}{
		Message: "password updated successfully",
	}

	err = json.NewEncoder(w).Encode(responce)

	if err != nil {
		myErr := utils.ErrorHandler(err, "error encoding json")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

}

func ResetPassExecHandler(w http.ResponseWriter, r *http.Request) {

}
