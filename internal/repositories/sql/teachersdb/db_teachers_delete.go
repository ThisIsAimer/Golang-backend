package teacherdb

import (
	"fmt"
	"net/http"
	"os"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"
	"strconv"
)

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
