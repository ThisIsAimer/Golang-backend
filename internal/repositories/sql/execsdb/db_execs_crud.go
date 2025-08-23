package execsdb

import (
	"net/http"
	"os"
	"reflect"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

func GetExecDBHandler(id int) (models.BasicExecs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.BasicExecs{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var executive models.BasicExecs
	executive.ID = id

	err = db.QueryRow("SELECT first_name, last_name, email, user_name, user_created_at, role, inactive_status class FROM execs WHERE id = ?", id).
		Scan(&executive.FirstName, &executive.LastName, &executive.Email, &executive.UserName, &executive.UserCreatedAt, &executive.Role, &executive.InactiveStatus)
	if err != nil {
		return models.BasicExecs{}, utils.ErrorHandler(err, "error retrieving data from database")
	}

	return executive, nil

}

func GetExecsDBHandler(r *http.Request, params []string) ([]models.BasicExecs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, user_name, user_created_at, role, inactive_status class FROM execs WHERE 1=1 `

	query, args := addFilters(r, query, params)

	query = addSorting(r, query, params)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error making query")
	}
	defer rows.Close()

	execs := make([]models.BasicExecs, 0)

	for rows.Next() {
		var exec models.BasicExecs
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.UserName, &exec.UserCreatedAt, &exec.Role, &exec.InactiveStatus)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning database results")
		}
		execs = append(execs, exec)
	}

	return execs, nil
}

// post-----------------------------------------------------------------------------------------------------------------------
func PostExecsDBHandler(modleTags []string, entries []models.Execs) ([]models.Execs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `INSERT INTO execs(`

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

		execsValue := reflect.ValueOf(&student).Elem()
		studentType := execsValue.Type()
		for i := range studentType.NumField() {
			if studentType.Field(i).Tag.Get("json") == "id,omitempty" {
				continue
			}
			arguments = append(arguments, execsValue.Field(i).Interface())
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

// patch-------------------------------------------------------------------------------------------
func PatchExecDBHandler(id int, arguments map[string]any) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "UPDATE execs SET "

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
