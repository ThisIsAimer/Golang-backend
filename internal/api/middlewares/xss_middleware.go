package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"simpleapi/pkg/utils"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// sanitize path-----------------------------------------------------------------------------------------
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}
		fmt.Println(sanitizedPath)

		// sanitize quary params-------------------------------------------------------------
		params := r.URL.Query()

		sanitizedQuery := make(map[string][]string)

		for k, values := range params {
			sanitizedKey, err := clean(k)
			if err != nil {
				http.Error(w, "query key is invalid", http.StatusInternalServerError)
				return
			}

			var sanatizedValues []string
			for _, v := range values {
				cleanValue, err := clean(v)
				if err != nil {
					http.Error(w, "query value is invalid", http.StatusInternalServerError)
					return
				}
				sanatizedValues = append(sanatizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanatizedValues
		}

		r.URL.Path = sanitizedPath.(string)

		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		fmt.Println("updated url:", r.URL.Path)
		fmt.Println("updated Query:", r.URL.RawQuery)

		next.ServeHTTP(w, r)
	})

}

//clean sanitizes input data

func clean(data any) (any, error) {

	switch d := data.(type) {
	case map[string]any:
		for k, v := range d {
			d[k] = sanitizeValue(v)
		}

		return d, nil

	case []any:
		for i, v := range d {
			d[i] = sanitizeValue(v)
		}

		return d, nil

	case string:
		return sanitizeString(d), nil

	default:
		return nil, utils.ErrorHandler(fmt.Errorf("unsupported type: %T", data), fmt.Sprintf("unsupported type: %T", data))
	}

}

func sanitizeValue(value any) any {

	switch d := value.(type) {

	case map[string]any:
		for k, v := range d {
			d[k] = sanitizeValue(v)
		}
		return d
	case []any:
		for i, v := range d {
			d[i] = sanitizeValue(v)
		}
		return d

	case string:
		return sanitizeString(d)

	default:
	}

	return 0
}

func sanitizeString(value string) string {

	return bluemonday.UGCPolicy().Sanitize(value)
}
