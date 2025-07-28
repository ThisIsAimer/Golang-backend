package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HppOptions struct {
	CheckQuary              bool
	CheckBody               bool
	CheckBodyForContentType string
	WhiteList               []string
}

func Hpp(options HppOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyForContentType) {
				filterBodyParams(r, options.WhiteList)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0]) //first value
			// r.Form.Set(k,v[len(v)-1]) //last value
		}

		if !isWhiteListerd(k, whitelist) {
			delete(r.Form, k)
		}
	}
}

func isWhiteListerd(param string, whitelist []string) bool {
	for _, v := range whitelist {
		if param == strings.ToLower(v) {
			return true
		}
	}
	return false
}
