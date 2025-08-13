package studentdb

import (
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

func GetStudentDBHandler(id int) (models.Student,error){
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var student models.Student
	student.ID = id

	err = db.QueryRow("SELECT first_name, last_name, email, class FROM students WHERE id = ?", id).
	Scan(&student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err,"error retrieving data from database")
	}


	return student, nil

}