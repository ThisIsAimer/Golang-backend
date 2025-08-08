package teacherdb

import (
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
)

func PostTeachersDBHandler(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return []models.Teacher{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?, ?, ?, ?, ?)")

	if err != nil {
		http.Error(w, "error in praparing sql query", http.StatusInternalServerError)
		return []models.Teacher{}, err
	}
	defer stmt.Close()

	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			http.Error(w, "error inserting values in the database (email may already exist)", http.StatusInternalServerError)
			return []models.Teacher{}, err
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "error getting last inserted id", http.StatusInternalServerError)
			return []models.Teacher{}, err
		}
		newTeachers[i].ID = int(lastId)

	}
	return newTeachers, nil
}
