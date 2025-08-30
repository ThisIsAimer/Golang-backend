package teachers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"simpleapi/pkg/utils"
)


func fieldIsEmpty(model any) error {
	element := reflect.ValueOf(model)
	for i := range element.NumField() {
		if element.Field(i).Kind() == reflect.String && element.Field(i).String() == "" {
			return utils.ErrorHandler(errors.New("user has not provided all fields"), "all fields required")
		}
	}

	return nil
}

func getPaginationParams(r *http.Request) (int, int) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	return page, limit
}