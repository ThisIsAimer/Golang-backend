package teachers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sqlconnect"
)

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {

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

	teacherVal := reflect.ValueOf(&existingTeacher).Elem() //the actual field of existingTeacher
	teacherType := teacherVal.Type()                       // type of the field(modles.Teacher)

	for k, v := range updates {
		for i := range teacherVal.NumField() {
			field := teacherType.Field(i) //returns each field of modles.Teacher with respect to i
			if field.Tag.Get("json") == k+`,omitempty` {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID,
	)

	if err != nil {
		http.Error(w, "error updating database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responce := struct {
		Status         string         `json:"status"`
		OldEntry       models.Teacher `json:"old_entry"`
		UpdatingValues map[string]any `json:"updated_values"`
		UpdatedEntry   models.Teacher `json:"updated_entries"`
	}{
		Status:         "success",
		OldEntry:       oldTeacher,
		UpdatingValues: updates,
		UpdatedEntry:   existingTeacher,
	}

	json.NewEncoder(w).Encode(responce)

}

// used for multi update
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var updates []map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		println("err is")
		http.Error(w, "error starting transaction", http.StatusInternalServerError)
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			http.Error(w, "error converting id to string", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			tx.Rollback()
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var existingTeacher models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)
		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "No rows found with the ID", http.StatusNotFound)
			return
		} else if err != nil {
			tx.Rollback()
			http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
			return
		}

		teacherVal := reflect.ValueOf(&existingTeacher).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := range teacherVal.NumField() {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+`,omitempty` {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(field.Type) {
							fieldVal.Set(val.Convert(field.Type))
						} else {
							tx.Rollback()
							http.Error(w, "cant convert value to value type", http.StatusInternalServerError)
						}
					} else {
						break
					}
				}
			}
			_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
				existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID,
			)
			if err != nil {
				http.Error(w, "error updating teacher", http.StatusInternalServerError)
				return
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "error commiting transaction to the client", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
