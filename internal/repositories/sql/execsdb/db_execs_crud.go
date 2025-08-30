package execsdb

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sql/sqlconnect"
	"simpleapi/pkg/utils"

	"github.com/go-mail/mail/v2"
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

func GetExecsDBHandler(r *http.Request, params []string, page, limit int) ([]models.BasicExecs, int, error) {
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, user_name, user_created_at, role, inactive_status class FROM execs WHERE 1=1 `

	query, args := addFilters(r, query, params)

	// pagination
	offset := (page - 1) * limit
	query += "LIMIT ? OFFSET ? "

	args = append(args, limit, offset)

	query = addSorting(r, query, params)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error making query")
	}
	defer rows.Close()

	execs := make([]models.BasicExecs, 0)

	for rows.Next() {
		var exec models.BasicExecs
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.UserName, &exec.UserCreatedAt, &exec.Role, &exec.InactiveStatus)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "error scanning database results")
		}
		execs = append(execs, exec)
	}

	// total count
	var totalCount int

	countQuery := `SELECT COUNT(*) FROM execs WHERE 1=1 `

	countQuery, newArgs := addFilters(r, countQuery, params)

	err = db.QueryRow(countQuery, newArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error getting total count")
	}

	return execs, totalCount, nil
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
			return nil, utils.ErrorHandler(err, "error making salt")
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
				return utils.ErrorHandler(err, "error making salt")
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

			if k == "user_name" {
				return utils.ErrorHandler(fmt.Errorf("trying to patch username"), ("cant edit user's username"))
			}

			if k == "password" {
				// hashing pass
				salt := make([]byte, 16)

				_, err := rand.Read(salt)
				if err != nil {
					return utils.ErrorHandler(err, "error making salt")
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
		return utils.ErrorHandler(err, "error making salt")
	}

	newHashedPassword, err := utils.PassEncoder(newPassword, salt)

	if err != nil {
		return utils.ErrorHandler(err, "error encoding new password")
	}

	query := `UPDATE execs SET password = ?, password_changed_at = ? WHERE id = ?`

	currentTime := time.Now().Format(time.RFC3339)

	_, err = db.Exec(query, newHashedPassword, currentTime, id)

	if err != nil {
		return utils.ErrorHandler(err, "error updating database")
	}

	return nil
}

// password reset Functions-------------------------------------------------------------------------------------
func ForgotPasswordDBHandler(email string) error {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var exec models.Execs

	err = db.QueryRow("SELECT id FROM execs WHERE email = ?", email).Scan(&exec.ID)

	if err != nil {
		return utils.ErrorHandler(err, "user not found")
	}

	expResetTime, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		return utils.ErrorHandler(err, "failed to send password reset mail")
	}

	mins := time.Duration(expResetTime) * time.Minute

	expiry := time.Now().Add(mins).Format(time.RFC3339)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return utils.ErrorHandler(err, "error making salt")
	}

	token := hex.EncodeToString(tokenBytes)

	hashedToken := sha256.Sum256(tokenBytes)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	_, err = db.Exec("UPDATE execs SET pass_reset_code = ?, pass_code_expires = ? WHERE id = ?", hashedTokenString, expiry, exec.ID)

	if err != nil {
		return utils.ErrorHandler(err, "error setting token")
	}

	resetUrl := fmt.Sprintf("https://localhost:3000/execs/login/resetpassword/reset/%s", token)
	message := fmt.Sprintf(" forgot your password? reset it using link %s \nIf you didnt reset a password reset, please ignore, the link is only valid for %v mins", resetUrl, expiry)

	myMail := mail.NewMessage()

	myMail.SetHeader("From", "schooladmin@school.com") // replace email
	myMail.SetHeader("To", email)
	myMail.SetHeader("Subject", "Password reset link")
	myMail.SetBody("text/plain", message)

	dialer := mail.NewDialer("localhost", 1025, "", "")
	err = dialer.DialAndSend(myMail)
	if err != nil {
		return utils.ErrorHandler(err, "error sending mail")
	}

	return nil
}

func ResetPassExecDBHandler(resetCode, new_pass string) error {

	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var exec models.Execs

	bytes, err := hex.DecodeString(resetCode)
	if err != nil {
		return utils.ErrorHandler(err, "error decoding string")
	}

	hashedToken := sha256.Sum256(bytes)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	query := `Select id, email FROM execs WHERE pass_reset_code = ? AND pass_code_expires > ?`

	err = db.QueryRow(query, hashedTokenString, time.Now().Format(time.RFC3339)).
		Scan(&exec.ID, &exec.Email)

	if err != nil {
		return utils.ErrorHandler(err, "invalid or expired reset code")
	}

	salt := make([]byte, 16)

	_, err = rand.Read(salt)
	if err != nil {
		return utils.ErrorHandler(err, "error making salt")
	}
	new_pass, err = utils.PassEncoder(new_pass, salt)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE execs SET password = ?, pass_reset_code = NULL, pass_code_expires = NULL, password_changed_at = ? where id = ?`

	_, err = db.Exec(updateQuery, new_pass, time.Now().Format(time.RFC3339), exec.ID)

	if err != nil {
		return utils.ErrorHandler(err, "error updating password")
	}

	return nil
}
