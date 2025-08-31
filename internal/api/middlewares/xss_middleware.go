package middlewares

import "net/http"

func XSSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}


//clean sanitizes input data

func clean(data any) (any, error ){

	switch data.(type){
	case map[string]any:

	case []any:
		
	case string:

	default:
	}

	return 0, nil
}