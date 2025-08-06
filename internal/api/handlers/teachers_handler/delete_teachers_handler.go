package teachers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"simpleapi/internal/repositories/sqlconnect"
)

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

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
