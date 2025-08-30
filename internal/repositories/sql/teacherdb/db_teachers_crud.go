package teacherdb

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
)

// get-----------------------------------------------------------------------------------------------------
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

func GetTeachersDBHandler(w http.ResponseWriter, r *http.Request, page, limit int) ([]models.Teacher, int, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error connecting to server")
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

	// pagination
	offset := (page - 1) * limit
	query += "LIMIT ? OFFSET ? "

	args = append(args, limit, offset)

	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error getting rows")
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "error scanning database results")
		}
		teacherList = append(teacherList, teacher)
	}

	// total count
	var totalCount int

	countQuery := `SELECT COUNT(*) FROM teachers WHERE 1=1 `

	countQuery, newArgs := addFilters(r, countQuery, queryParams)

	err = db.QueryRow(countQuery, newArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error getting total count")
	}

	return teacherList, totalCount, nil
}

// post------------------------------------------------------------------------------------------------------------------------------------
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

// put-----------------------------------------------------------------------------------------------------------------------------------------------------------
func PutTeacherDBHandler(w http.ResponseWriter, id int, updatedTeacher models.Teacher) (models.Teacher, models.Teacher, error) {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return models.Teacher{}, models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	existingTeacher, err := GetExistingTeacher(w, db, id)
	if err != nil {
		return models.Teacher{}, models.Teacher{}, err
	}
	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID,
	)

	if err != nil {
		return models.Teacher{}, models.Teacher{}, utils.ErrorHandler(err, "error updating database")
	}

	return updatedTeacher, existingTeacher, nil

}

// patch-----------------------------------------------------------------------------------------------
func PatchTeacherDBHandler(w http.ResponseWriter, id int, updates map[string]any) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)

	if err == sql.ErrNoRows {
		return utils.ErrorHandler(err, "no rows found")
	} else if err != nil {
		return utils.ErrorHandler(err, "Unable to retrieve data")
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem() //the actual field of existingTeacher
	teacherType := teacherVal.Type()                       // type of the field(modles.Teacher)

	for k, v := range updates {
		for i := range teacherVal.NumField() {
			field := teacherType.Field(i) //returns each field of modles.Teacher with respect to i
			if field.Tag.Get("json") == k+`,omitempty` {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID,
	)

	if err != nil {
		return utils.ErrorHandler(err, "error updating database")
	}

	return nil

}

func PatchTeachersDBHandler(w http.ResponseWriter, updates []map[string]any) error {

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

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid id")
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "invalid ID")
		}

		var existingTeacher models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(
			&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject,
		)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.ErrorHandler(err, "no rows found")
		} else if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "unable to retrieve data")
		}

		teacherVal := reflect.ValueOf(&existingTeacher).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := range teacherVal.NumField() {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+`,omitempty` {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(field.Type) {
							fieldVal.Set(val.Convert(field.Type))
						} else {
							tx.Rollback()
							return utils.ErrorHandler(err, "unconvertable type")
						}
					} else {
						break
					}
				}
			}
			_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
				existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID,
			)
			if err != nil {
				tx.Rollback()
				return utils.ErrorHandler(err, "error updating database")
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error commiting transaction")
	}

	return nil
}

// delete---------------------------------------------------------------------------------------------
func DeleteTeacherDBHandler(w http.ResponseWriter, id int) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to server")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "error deleting row")
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error retrieveing deleted results")
	}
	if rowsEffected == 0 {
		return utils.ErrorHandler(err, "row now found")
	}

	return nil

}

func DeleteTeachersDBHandler(w http.ResponseWriter, ids []string) ([]int, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		return nil, utils.ErrorHandler(err, "error startingTransaction")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		return nil, utils.ErrorHandler(err, "error preparing delete statement")
	}
	defer stmt.Close()

	var deletedIds []int

	for _, value := range ids {
		id, err := strconv.Atoi(value)
		if err != nil {
			return nil, utils.ErrorHandler(err, "invalid ID")
		}

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error executing statement")
		}

		deletedRows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error retrieveing delete results")
		}

		if deletedRows > 0 {
			deletedIds = append(deletedIds, id)
		}
		if deletedRows == 0 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("%d id doesnt exist", id))
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error commiting transaction")
	}

	return deletedIds, nil
}

// students assigned to teacher --------------------------------------------------------------------------------
func GetStudentsByTeacherIdDB(id string) (models.Teacher, []models.Student, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Teacher{}, nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var teacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, nil, utils.ErrorHandler(err, "invalid ID")
	} else if err != nil {
		return models.Teacher{}, nil, utils.ErrorHandler(err, "database error")
	}

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE class = (SELECT class FROM teachers where id = ?)"

	rows, err := db.Query(query, id)
	if err != nil {
		return models.Teacher{}, nil, utils.ErrorHandler(err, "error getting rows")
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return models.Teacher{}, nil, utils.ErrorHandler(err, "error scanning database results")
		}
		students = append(students, student)
	}

	return teacher, students, nil
}

// students count------------------------------------------------------------------------------------------
func GetStudentCountByTeacherIdDB(id string) (int, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return 0, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	//agrigate functions in sql!
	query := "SELECT COUNT(*) FROM students WHERE class = (SELECT class FROM teachers where id = ?)"
	var studentCount int

	err = db.QueryRow(query, id).Scan(&studentCount)
	if err != nil {
		return 0, utils.ErrorHandler(err, "error getting rows")
	}

	return studentCount, nil
}
