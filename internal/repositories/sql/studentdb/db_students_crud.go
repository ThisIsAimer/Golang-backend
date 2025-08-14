package studentdb

import (
	"net/http"
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

func GetStudentsDBHandler(r *http.Request, params []string) ([]models.Student, error){
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, class FROM students WHERE 1 = 1 `

	query, args := addFilters(r,query,params)

	query = addSorting(r,query,params)

	rows, err := db.Query(query,args...)
	if err != nil {
		return nil , utils.ErrorHandler(err, "error making query")
	}

	students := make([]models.Student,0)

	for rows.Next(){
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning database results")
		}
		students = append(students, student)
	}

	return students, nil
}