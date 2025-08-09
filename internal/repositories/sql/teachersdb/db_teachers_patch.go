package teacherdb

import (
	"database/sql"
	"net/http"
	"os"
	"reflect"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
	"strconv"
)

func PatchTeacherDBHandler(w http.ResponseWriter, id int, updates map[string]any) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	if err == sql.ErrNoRows {
		return utils.ErrorHandler(err, "no rows found")
	} else if err != nil {
		return utils.ErrorHandler(err, "Unable to retrieve data")
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
		return utils.ErrorHandler(err, "error updating database")
	}

	return nil

}

func PatchTeachersDBHandler(w http.ResponseWriter, updates []map[string]any) error {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "error starting transaction")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid id")
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid ID")
		}

		var existingTeacher models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.ErrorHandler(err, "no rows found")
		} else if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "unable to retrieve data")
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
							return utils.ErrorHandler(err, "unconvertable type")
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
				return utils.ErrorHandler(err, "error updating database")
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error commiting transaction")
	}

	return nil
}
