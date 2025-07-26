package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func CompMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request)  {
		// check if the client accepts gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"),"gzip"){
			next.ServeHTTP(w,r)
		} else{
			w.Header().Set("Content-Encoding","gzip")
			w.Header().Del("Content-Length") // content length might cause some problems

			
			fmt.Println("compression active")


			gz := gzip.NewWriter(w)
			defer gz.Close()

			// wrap the responsewriter
			gzipRW := &gzipResponseWriter{ResponseWriter: w, Writer: gz}

			next.ServeHTTP(gzipRW,r)

		}

	})

}

// gzip response writer will make ResponseWriter write gunzip responses

type gzipResponseWriter struct{
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gz *gzipResponseWriter) Write(b []byte) (int,error){
	return gz.Writer.Write(b)
}