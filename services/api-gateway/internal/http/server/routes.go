package server

import (
	"net/http"

	middleware "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []middleware.Middleware
}
