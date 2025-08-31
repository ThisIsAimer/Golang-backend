package middlewares

import (
	"fmt"
	"net/http"
	"simpleapi/pkg/utils"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// sanitize path-----------------------------------------------------------------------------------------
		sanitizePath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid path", http.StatusInternalServerError)
			return
		}
		fmt.Println(sanitizePath)

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
