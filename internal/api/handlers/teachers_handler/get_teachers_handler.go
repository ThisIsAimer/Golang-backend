package teachers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sqlconnect"
	"strconv"
	"strings"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	

	w.Header().Set("Content-Type", "application/json")

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1 "

	queryParams := []string{
		"first_name",
		"last_name",
		"email",
		"class",
		"subject",
	}

	query, args := addFilters(r, query, queryParams)

	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "error getting rows", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "error scanning database results", http.StatusInternalServerError)
			return
		}
		teacherList = append(teacherList, teacher)
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

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idstr := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")

	//handle path parametre
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	
	}

}

func addFilters(r *http.Request, query string, params []string) (string, []any) {

	var args []any

	for _, value := range params {
		result := r.URL.Query().Get(value)
		if result != "" {
			query += "AND " + value + "= ? "
			args = append(args, result)
		}
	}

	return query, args

}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) != 0 {
		for i, params := range sortParams {
			parts := strings.Split(params, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidsortOrder(order) || !isValidField(field) {
				continue
			}

			query += "ORDER BY "

			if i > 0 {
				query += " , "
			}

			query += field + " " + order

		}
	}
	return query
}

func isValidsortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidField(field string) bool {
	validFields := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"class":      true,
		"email":      true,
		"subject":    true,
	}
	return validFields[field]
}
