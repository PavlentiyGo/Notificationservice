package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	grpc_client "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/grpc"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/handlers"
	middleware "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/server"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.NewConfigMust()

	serv := server.NewServer(cfg)

	subscriptionConn, err := grpc_client.NewGRPCConn(cfg.SubscriptionAddr)
	if err != nil {
		log.Fatalf("failed to start subs client: %s", err)
	}
	defer subscriptionConn.Close()

	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionConn, cfg)
	serv.ChainRoutes(subscriptionHandler.Routes()...)

	analysisConn, err := grpc_client.NewGRPCConn(cfg.AnalysisAddr)
	if err != nil {
		log.Fatalf("failed to start subs client: %s", err)
	}
	defer analysisConn.Close()

	analysisHandler := handlers.NewAnalysisHandler(cfg, analysisConn)
	serv.ChainRoutes(analysisHandler.Routes()...)

	if err = serv.Run(
		ctx,
		middleware.CORS([]string{
			"http://localhost:5173",
			"https://frontsubscriptionreminder.vercel.app/",
			"https://frontsubscriptionreminder.vercel.app",
		}),
		middleware.Trace(),
		middleware.Recover(),
	); err != nil {
		log.Fatal("failed to run server", err)
	}
}
