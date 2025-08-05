package teachers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sqlconnect"
)

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idstr := strings.TrimSuffix(path, "/")

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
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	oldTeacher := existingTeacher

	if err == sql.ErrNoRows {
		http.Error(w, "No rows found with the ID", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}
	query := "UPDATE teachers SET "

	for k, v := range updates {
		switch k {
		case "first_name":
			existingTeacher.FirstName = v.(string)
			query += "first_name, "
		case "last_name":
			existingTeacher.LastName = v.(string)
		case "email":
			existingTeacher.Email = v.(string)
		case "class":
			existingTeacher.Class = v.(string)
		case "subject":
			existingTeacher.Subject = v.(string)
		}
	}

	//applying updates using reflect
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	fmt.Println("teacher:", teacherVal.Type())

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID,
	)

	if err != nil {
		http.Error(w, "error updating database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Status         string `json:"status"`
		OldEntry       models.Teacher
		UpdatingValues map[string]any
		UpdatedEntry   models.Teacher
	}{
		Status:         "success",
		OldEntry:       oldTeacher,
		UpdatingValues: updates,
		UpdatedEntry:   existingTeacher,
	}

	json.NewEncoder(w).Encode(responce)

}
