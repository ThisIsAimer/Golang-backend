package execsdb

import (
	"crypto/rand"
	"database/sql"
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

		// password hashing
		salt := make([]byte, 16)

		_, err := rand.Read(salt)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		student.Password, err = utils.PassEncoder(student.Password, salt)
		if err != nil {
			return nil, err
		}

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

		if k == "password" {

			// hashing pass
			salt := make([]byte, 16)

			_, err := rand.Read(salt)
			if err != nil {
				return utils.ErrorHandler(err, "error adding data")
			}

			v, err = utils.PassEncoder(v.(string), salt)
			if err != nil {
				return err
			}
		}

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

func PatchExecsDBHandler(argumentsList []map[string]any) error {
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
		query := "UPDATE execs SET "

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

			if k == "password" {
				// hashing pass
				salt := make([]byte, 16)

				_, err := rand.Read(salt)
				if err != nil {
					return utils.ErrorHandler(err, "error adding data")
				}
				v, err = utils.PassEncoder(v.(string), salt)
				if err != nil {
					return err
				}
			}

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

// Delete----------------------------------------------------------------------------------------
func DeleteExecDBHandler(id int) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
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

func DeleteExecsDBHandler(ids []int) error {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "DELETE FROM execs WHERE id IN "

	args := "("

	var anyIds []any

	for _, id := range ids {
		if args != "(" {
			args += ", "
		}
		args += "?"
		anyIds = append(anyIds, id)
	}
	args += ")"

	query += args

	result, err := db.Exec(query, anyIds...)

	if err != nil {
		return utils.ErrorHandler(err, "error executing statement")
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error retrieveing delete results")
	}
	if deletedRows == 0 {
		return utils.ErrorHandler(err, " one of the ids doesnt exist")
	}

	return nil
}

// Login-------------------------------------------------------------------------------------------------------
func LoginExecDBHandler(username string) (models.Execs, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return models.Execs{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var user models.Execs

	err = db.QueryRow("SELECT id, user_name, password, first_name, last_name, email, user_created_at, role, inactive_status class FROM execs WHERE user_name = ?", username).
		Scan(&user.ID, &user.UserName, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.UserCreatedAt, &user.Role, &user.InactiveStatus)

	if err == sql.ErrNoRows {
		return models.Execs{}, utils.ErrorHandler(err, "invalid username")
	} else if err != nil {
		return models.Execs{}, utils.ErrorHandler(err, "error getting data")
	}

	return user, nil
}

// passwords

func UpdatePassExecDBHandler(id int, currentPassword, newPassword string) error {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var existingUsername string
	var existingPass string

	err = db.QueryRow("SELECT user_name, password FROM execs WHERE id = ?", id).
		Scan(&existingUsername, &existingPass)

	if err == sql.ErrNoRows {
		return utils.ErrorHandler(err, "invalid user id")
	} else if err != nil {
		return utils.ErrorHandler(err, "error scanning execs")
	}

	err = utils.VerifyPassword(currentPassword, existingPass)

	if err != nil {
		return utils.ErrorHandler(err, "incorrect current password")
	}

	salt := make([]byte, 16)

	_, err = rand.Read(salt)
	if err != nil {
		return utils.ErrorHandler(err, "error adding data")
	}

	newHashedPassword, err := utils.PassEncoder(newPassword, salt)

	if err != nil {
		return utils.ErrorHandler(err, "error encoding new password")
	}

	query := `UPDATE execs SET password = ? WHERE id = ?`

	_, err = db.Exec(query, newHashedPassword, id)

	if err != nil {
		return utils.ErrorHandler(err, "error updating database")
	}

	return nil
}
