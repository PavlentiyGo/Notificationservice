package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
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

func (s *Server) ChainRoutes(routes ...Route) {
	for _, route := range routes {
		prefix := fmt.Sprintf("%s %s", route.Method, route.Path)
		handler := middleware.ChainMiddlewares(route.Handler, route.Middlewares...)
		s.Handle(prefix, handler)
	}
}

func (s *Server) Run(
	ctx context.Context,
	middlewares ...middleware.Middleware,
) error {

	mux := middleware.ChainMiddlewares(s.ServeMux, middlewares...)

	server := &http.Server{
		Handler: mux,
		Addr:    s.config.ServerAddr,
	}

	ch := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("failed to listen and server server: %s", err)
			ch <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(ctx, s.config.GracefulTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown server properly: %s", err)
			server.Close()
		}
	case err := <-ch:
		return err
	}
	return nil
}
