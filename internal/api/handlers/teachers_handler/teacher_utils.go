package teachers

import (
	"errors"
	"reflect"
	
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