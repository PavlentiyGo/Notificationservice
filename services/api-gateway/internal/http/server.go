package router

import (
	"net/http"

	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
)

type Server struct {
	*http.ServeMux
	config config.Config
}

func NewServer(config config.Config) *Server {
	return &Server{
		ServeMux: http.NewServeMux(),
		config:   config,
	}
}

func ChainRoutes() {}
