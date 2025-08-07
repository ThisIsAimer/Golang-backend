package teachers

import (
	"encoding/json"
	"net/http"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/teachersdb"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	teacherList, err := teacher_db.GetTeachersDBHandler(w, r)
	if err != nil {
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Content-Type", "application/json")
	idstr := r.PathValue("id")

	//handle path parametre
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	teacher, err := teacher_db.GetTeacherDBHandler(w, r, id)

	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return

	}

}
