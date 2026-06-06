package middleware

import (
	"log"
	"net/http"
	"time"
)

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			timeStart := time.Now()

			log.Printf("got request: Path:%s, Method:%s, TimeStart:%s", r.URL.Path, r.Method, timeStart)
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     200,
			}
			next.ServeHTTP(rw, r)
			log.Printf("end request: Path:%s, Method:%s, SpendTime:%s, StatusCode:%d", r.URL.Path, r.Method, time.Since(timeStart), rw.statusCode)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code

}
