package middleware

import (
	"fmt"
	"net/http"

	http2 "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/response"
)

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				p := recover()
				responseHandler := http2.NewResponseHandler(w)
				if p != nil {
					responseHandler.ErrorResponse(fmt.Sprintf("got unexpected panic: %s", p), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
