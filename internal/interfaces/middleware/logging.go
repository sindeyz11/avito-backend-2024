package middleware

import (
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		lrw := negroni.NewResponseWriter(w)

		next.ServeHTTP(lrw, request)

		statusCode := lrw.Status()

		log.Printf("%d %s %s %s", statusCode, request.Method, request.RequestURI, time.Since(start))
	})
}
