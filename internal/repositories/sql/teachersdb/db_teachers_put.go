package teacherdb

import (
	"database/sql"
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
)

func PutTeacherDBHandler(w http.ResponseWriter, id int, updatedTeacher models.Teacher) (models.Teacher, models.Teacher, error) {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return models.Teacher{}, models.Teacher{}, err
	}
	defer db.Close()

	existingTeacher, err := GetExistingTeacher(w, db, id)
	if err != nil {
		return models.Teacher{},models.Teacher{}, err
	}
	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID,
	)

	if err != nil {
		http.Error(w, "error updating database", http.StatusInternalServerError)
		return models.Teacher{},models.Teacher{}, err
	}
	return updatedTeacher, existingTeacher, err

}

func GetExistingTeacher(w http.ResponseWriter, db *sql.DB, id int) (models.Teacher, error) {
	var existingTeacher models.Teacher

	err := db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	if err == sql.ErrNoRows {
		http.Error(w, "No rows found with the ID", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return existingTeacher, nil
}
