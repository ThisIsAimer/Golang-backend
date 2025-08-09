package teacherdb

import (
	"database/sql"
	"net/http"
	"os"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
	"strings"
)

func GetTeachersDBHandler(w http.ResponseWriter, r *http.Request) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to server")
	}
	defer db.Close()
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1 "

	queryParams := []string{
		"first_name",
		"last_name",
		"email",
		"class",
		"subject",
	}

	query, args := addFilters(r, query, queryParams)

	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error getting rows")
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning database results")
		}
		teacherList = append(teacherList, teacher)
	}
	return teacherList, nil
}

func addFilters(r *http.Request, query string, params []string) (string, []any) {

	var args []any

	for _, value := range params {
		result := r.URL.Query().Get(value)
		if result != "" {
			query += "AND " + value + "= ? "
			args = append(args, result)
		}
	}

	return query, args

}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) != 0 {
		for i, params := range sortParams {
			parts := strings.Split(params, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidsortOrder(order) || !isValidField(field) {
				continue
			}

			query += "ORDER BY "

			if i > 0 {
				query += " , "
			}

			query += field + " " + order

		}
	}
	return query
}

func isValidsortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidField(field string) bool {
	validFields := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"class":      true,
		"email":      true,
		"subject":    true,
	}
	return validFields[field]
}

func GetTeacherDBHandler(w http.ResponseWriter, id int) (models.Teacher, error) {
	var teacher models.Teacher

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "invalid ID")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "database error")
	}

	return teacher, nil

}
