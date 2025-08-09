package teachers

import (
	"encoding/json"
	"net/http"
	"os"
	"simpleapi/internal/repositories/sql/sqlconnect"
	teacherdb "simpleapi/internal/repositories/sql/teachersdb"
	"strconv"
)

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

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
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var ids []string

	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
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
