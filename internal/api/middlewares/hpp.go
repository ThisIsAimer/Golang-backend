package middlewares

import "net/http"

type HppOptions struct{
	CheckQuary bool
	CheckBody bool
	CheckBodyForContentType string
	WhiteList []string

}


func Hpp(options HppOptions) func(http.Handler) http.Handler{
	return func (next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}