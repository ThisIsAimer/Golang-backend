package execsdb

import (
	"net/http"
	"os"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

func GetExecDBHandler(id int) (models.Execs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Execs{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var executive models.Execs
	executive.ID = id

	err = db.QueryRow("SELECT first_name, last_name, email, user_name, user_created_at, inactive_status class FROM execs WHERE id = ?", id).
		Scan(&executive.FirstName, &executive.LastName, &executive.Email, &executive.UserName, &executive.UserCreatedAt, &executive.InactiveStatus)
	if err != nil {
		return models.Execs{}, utils.ErrorHandler(err, "error retrieving data from database")
	}

	return executive, nil

}

func GetExecsDBHandler(r *http.Request, params []string) ([]models.Execs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, user_name, user_created_at, inactive_status class FROM execs WHERE 1=1 `

	query, args := addFilters(r, query, params)

	query = addSorting(r, query, params)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error making query")
	}
	defer rows.Close()

	execs := make([]models.Execs, 0)

	for rows.Next() {
		var exec models.Execs
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.UserName, &exec.UserCreatedAt, &exec.InactiveStatus)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning database results")
		}
		execs = append(execs, exec)
	}

	return execs, nil
}
