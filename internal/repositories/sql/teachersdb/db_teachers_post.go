package teacherdb

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
	"strings"
)

func PostTeachersDBHandler(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return []models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	// 	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?, ?, ?, ?, ?)")
	stmt, err := db.Prepare(getInsertQuery(models.Teacher{}))

	if err != nil {
		return []models.Teacher{}, utils.ErrorHandler(err, "error preparing statement")
	}
	defer stmt.Close()

	for i, teacher := range newTeachers {
		// res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		res, err := stmt.Exec(getExecValues(teacher)...)
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
