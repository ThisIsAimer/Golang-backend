package execsdb

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"simpleapi/pkg/utils"
	"strings"

	"golang.org/x/crypto/argon2"
)

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

func addSorting(r *http.Request, query string, validFields []string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) != 0 {
		for i, params := range sortParams {
			parts := strings.Split(params, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidsortOrder(order) || !isValidField(validFields, field) {
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

func isValidField(validfields []string, field string) bool {
	boolFields := map[string]bool{}

	for _, v := range validfields {
		boolFields[v] = true
	}

	return boolFields[field]
}

func passEncoder(password string) (string, error) {
	if password == "" {
		return "", utils.ErrorHandler(fmt.Errorf("password is empty"), "password is required")
	}

	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return "", utils.ErrorHandler(err, "error adding data")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := saltBase64 + "." + hashBase64

	return encodedHash, nil
}
