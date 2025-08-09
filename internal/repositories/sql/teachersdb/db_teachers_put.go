package teacherdb

import (
	"database/sql"
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

func PutTeacherDBHandler(w http.ResponseWriter, id int, updatedTeacher models.Teacher) (models.Teacher, models.Teacher, error) {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return models.Teacher{}, models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	existingTeacher, err := GetExistingTeacher(w, db, id)
	if err != nil {
		return models.Teacher{}, models.Teacher{}, err
	}
	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID,
	)

	if err != nil {
		return models.Teacher{}, models.Teacher{}, utils.ErrorHandler(err, "error updating database")
	}

	return updatedTeacher, existingTeacher, nil

}

func GetExistingTeacher(w http.ResponseWriter, db *sql.DB, id int) (models.Teacher, error) {
	var existingTeacher models.Teacher

	err := db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "rows not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "unable to retrieve data")
	}
	return existingTeacher, nil
}
