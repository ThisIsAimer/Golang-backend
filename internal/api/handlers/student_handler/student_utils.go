package students

import (
	"reflect"
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
