package teachers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"simpleapi/internal/repositories/sqlconnect"
)

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "error deleting row", http.StatusInternalServerError)
		return
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "error retrieve delete result", http.StatusInternalServerError)
		return
	}
	if rowsEffected == 0 {
		http.Error(w, "row not found", http.StatusNotFound)
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

	tx, err := db.Begin()

	if err != nil {
		http.Error(w, "error starting transaction", http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		http.Error(w, "error prapring delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var deletedIds []int

	for _, value := range ids {
		id, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			http.Error(w, "error executing statement", http.StatusInternalServerError)
			return
		}

		deletedRows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			http.Error(w, "error retrieveing delete result", http.StatusInternalServerError)
			return
		}

		if deletedRows > 0 {
			deletedIds = append(deletedIds, id)
		}
		if deletedRows == 0 {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("%d id doesnt exist", id), http.StatusBadRequest)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "error commiting transaction", http.StatusInternalServerError)
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
