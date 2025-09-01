package teachers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/teacherdb"
	"simpleapi/pkg/utils"
)

// CRUD GET-----------------------------------------------------------------------------------------------
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	page, limit := getPaginationParams(r)

	teacherList, count, err := teacherdb.GetTeachersDBHandler(w, r, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startEntry := ((page - 1) * limit)
	endEntry := startEntry + limit

	startEntry++

	if endEntry > count {
		endEntry = count
	}
	if startEntry > count {
		startEntry = 0
		endEntry = 0
	}

	strCount := fmt.Sprintf("%d-%d of %d", startEntry, endEntry, count)

	response := struct {
		Status   string           `json:"status"`
		Count    string           `json:"count"`
		PageNo   int              `json:"page_no"`
		PageSize int              `json:"page_size"`
		Data     []models.Teacher `json:"data"`
	}{
		Status:   "success",
		Count:    strCount,
		PageNo:   page,
		PageSize: limit,
		Data:     teacherList,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	idstr := r.PathValue("id")

	//handle path parametre
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	teacher, err := teacherdb.GetTeacherDBHandler(w, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return

	}

}

// CRUD POST------------------------------------------------------------------------------------------------
func PostTeachersHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin", "moderator")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	// used to discard unknown fields
	decoder.DisallowUnknownFields()

	var newTeachers []models.Teacher
	err = decoder.Decode(&newTeachers)
	if err != nil {
		myError := utils.ErrorHandler(err, "invalid request body")
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

// CRUD PUT--------------------------------------------------------------------------------------
func PutTeacherHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin", "moderator")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

// CRUD PATCH-------------------------------------------------------------------------------------------------------
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin", "moderator")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	err = teacherdb.PatchTeacherDBHandler(w, id, updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// used for multi update
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin", "moderator")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var updates []map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request payload json not in format", http.StatusBadRequest)
		return
	}
	err = teacherdb.PatchTeachersDBHandler(w, updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// CRUD DELETE--------------------------------------------------------------------------------------------------
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	err = teacherdb.DeleteTeacherDBHandler(w, id)
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

	json.NewEncoder(w).Encode(responce)

}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var ids []string
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		myErr := utils.ErrorHandler(err, "error decoding ids")
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	deletedIds, err := teacherdb.DeleteTeachersDBHandler(w, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(deletedIds) == 0 {
		http.Error(w, "ids dont exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "success",
		DeletedIDs: deletedIds,
	}

	json.NewEncoder(w).Encode(response)
}

// students by teacher id

func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	teacher, students, err := teacherdb.GetStudentsByTeacherIdDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status       string           `json:"status"`
		Teacher      models.Teacher   `json:"teacher"`
		StudentCount int              `json:"student_count"`
		Students     []models.Student `json:"students"`
	}{
		Status:       "SUCCESS",
		Teacher:      teacher,
		StudentCount: len(students),
		Students:     students,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetStudentCountByTeacherId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	studentCount, err := teacherdb.GetStudentCountByTeacherIdDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status       string `json:"status"`
		StudentCount int    `json:"student_count"`
	}{
		Status:       "SUCCESS",
		StudentCount: studentCount,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
