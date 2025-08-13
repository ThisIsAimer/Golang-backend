package teacherdb

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"simpleapi/internal/models"
	"simpleapi/pkg/utils"
	"strings"
)

// get-----------------------------------------------------------------------------------------------------------------------
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

// post
func getInsertQuery(model any) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholder string

	for i := range modelType.NumField() {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")

		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholder += ", "
			}
			columns += dbTag
			placeholder += "?"
		}
	}

	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES(%s)", columns, placeholder)
}

func getExecValues(modle any) []any {
	modleValue := reflect.ValueOf(modle)
	modleType := modleValue.Type()
	values := []any{}

	for i := range modleType.NumField() {
		dbTag := modleType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modleValue.Field(i).Interface())
		}
	}

	return values
}


// put---------------------------------------------------------------------------------------------------------
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
