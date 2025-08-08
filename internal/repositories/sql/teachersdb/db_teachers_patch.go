package teacherdb

import (
	"database/sql"
	"net/http"
	"os"
	"reflect"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"strconv"
)

func PatchTeacherDBHandler(w http.ResponseWriter, id int, updates map[string]any) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	if err == sql.ErrNoRows {
		http.Error(w, "No rows found with the ID", http.StatusNotFound)
		return err
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return err
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
		return err
	}

	return nil

}

func PatchTeachersDBHandler(w http.ResponseWriter, updates []map[string]any) error {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "error starting transaction", http.StatusInternalServerError)
		return err
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			http.Error(w, "error converting id to string", http.StatusBadRequest)
			return err
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			tx.Rollback()
			http.Error(w, "invalid id", http.StatusBadRequest)
			return err
		}

		var existingTeacher models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)
		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "No rows found with the ID", http.StatusNotFound)
			return err
		} else if err != nil {
			tx.Rollback()
			http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
			return err
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
				return err
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "error commiting transaction to the client", http.StatusInternalServerError)
		return err
	}

	return nil
}
