package studentdb

import (
	"net/http"
	"os"
	"reflect"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

// get-------------------------------------------------------------------------------------------------------
func GetStudentDBHandler(id int) (models.Student, error) {
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
		return models.Student{}, utils.ErrorHandler(err, "error retrieving data from database")
	}

	return student, nil

}

func GetStudentsDBHandler(r *http.Request, params []string) ([]models.Student, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, class FROM students WHERE 1 = 1 `

	query, args := addFilters(r, query, params)

	query = addSorting(r, query, params)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error making query")
	}

	students := make([]models.Student, 0)

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning database results")
		}
		students = append(students, student)
	}

	return students, nil
}

// post ----------------------------------------------------------------------------------------------------------------------
func PostStudentsDBHandler(modleTags []string, entries []models.Student) ([]models.Student, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `INSERT INTO students(`

	params := ""
	for _, v := range modleTags {
		if v == "id" {
			continue
		}
		if params != "" {
			params += ", "
		}
		params += v
	}

	query += params + ") VALUES"

	givenValues := ""
	var arguments []any

	for _, student := range entries {
		if givenValues != "" {
			givenValues += ", "
		}
		startBracket := "("
		for _, v := range modleTags {
			if v == "id" {
				continue
			}
			if startBracket != "(" {
				startBracket += ", "
			}
			startBracket += "?"
		}
		startBracket += ")"
		givenValues += startBracket

		studentValue := reflect.ValueOf(&student).Elem()
		studentType := studentValue.Type()
		for i := range studentType.NumField() {
			if studentType.Field(i).Tag.Get("json") == "id,omitempty" {
				continue
			}
			arguments = append(arguments, studentValue.Field(i).Interface())
		}

	}

	query += givenValues

	result, err := db.Exec(query, arguments...)
	if err != nil {
		myErr := utils.ErrorHandler(err, "error uploading entries")
		return nil, myErr
	}
	num, err := result.LastInsertId()
	id := int(num)

	if err != nil {
		myErr := utils.ErrorHandler(err, `error getting lastValue`)
		return nil, myErr
	}

	for i := range entries {
		entries[i].ID = id
		id++
	}

	return entries, nil
}

// put-----------------------------------------------------------------------------------
func PutStudentDBHandler(id int, entry models.Student, params []string) (models.Student, models.Student, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Student{}, models.Student{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "SELECT "

	for _, v := range params {
		if v == "id" {
			continue
		}
		if query != "SELECT " {
			query += ", "
		}
		query += v
	}
	query += " FROM students WHERE id = ?"

	var existingStudent models.Student

	err = db.QueryRow(query, id).
		Scan(&existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)

	if err != nil {
		return models.Student{}, models.Student{}, utils.ErrorHandler(err, "error retrieving data from database")
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		entry.FirstName, entry.LastName, entry.Email, entry.Class, id,
	)
	entry.ID = id

	if err != nil {
		return models.Student{}, models.Student{}, utils.ErrorHandler(err, "error updating database")
	}

	return entry, existingStudent, nil
}

// patch ---------------------------------------------------------------------------------

func PatchStudentDBHandler(id int, arguments map[string]any) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "UPDATE students SET "

	flags := ""
	var args []any

	for k, v := range arguments {
		if flags != "" {
			flags += ", "
		}
		flags += k + " = ?"

		args = append(args, v)

	}
	args = append(args, id)

	query += flags + " WHERE id = ?"

	_, err = db.Exec(query, args...)

	if err != nil {
		return utils.ErrorHandler(err, "error updating database")
	}

	return nil

}

func PatchStudentsDBHandler(argumentsList []map[string]any) error {
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

	for _, arguments := range argumentsList {
		query := "UPDATE students SET "

		flags := ""
		var args []any
		var id any

		for k, v := range arguments {
			if k == "id" {
				id = v
				continue
			}
			if flags != "" {
				flags += ", "
			}
			flags += k + " = ?"

			args = append(args, v)

		}
		args = append(args, id)

		query += flags + " WHERE id = ?"

		_, err = tx.Exec(query, args...)

		if err != nil {
			return utils.ErrorHandler(err, "error updating database")
		}

	}
	tx.Commit()

	return nil
}
