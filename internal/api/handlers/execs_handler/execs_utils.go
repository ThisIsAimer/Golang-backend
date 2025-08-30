package execs

import (
	"net/http"
	"reflect"
	"strconv"
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
