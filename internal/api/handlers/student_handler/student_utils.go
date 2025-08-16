package students

import (
	"errors"
	"reflect"
	"simpleapi/pkg/utils"
	"strings"
)

func getModelTags(model any) []string {
	tags := []string{}
	modelType := reflect.TypeOf(model)
	for i := range modelType.NumField() {
		tag := modelType.Field(i).Tag.Get("json")
		tag = strings.TrimSuffix(tag, `,omitempty`)
		tags = append(tags, tag)
	}
	
	return tags
}

func fieldIsEmpty(model any) error {
	element := reflect.ValueOf(model)
	for i := range element.NumField() {
		if element.Field(i).Kind() == reflect.String && element.Field(i).String() == "" {
			return utils.ErrorHandler(errors.New("user has not provided all fields"), "all fields required")
		}
	}

	return nil
}
