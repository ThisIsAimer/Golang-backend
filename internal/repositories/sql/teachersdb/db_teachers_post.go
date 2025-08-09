package teacherdb

import (
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

func PostTeachersDBHandler(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return []models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?, ?, ?, ?, ?)")

	if err != nil {
		return []models.Teacher{}, utils.ErrorHandler(err, "error preparing statement")
	}
	defer stmt.Close()

	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			return []models.Teacher{}, utils.ErrorHandler(err, "error posting data")
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return []models.Teacher{}, utils.ErrorHandler(err, "error getting last id")
		}
		newTeachers[i].ID = int(lastId)

	}
	return newTeachers, nil
}
