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

	subscriptionClient, err := grpc_client.NewGRPCConn(cfg.SubscriptionAddr)
	if err != nil {
		log.Fatalf("failed to start subs client: %s", err)
	}
	defer subscriptionClient.Close()

	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionClient, cfg)
	serv.ChainRoutes(subscriptionHandler.Routes()...)

	if err = serv.Run(
		ctx,
		middleware.Recover(),
		middleware.CORS([]string{
			"http://localhost:5173",
			"https://frontsubscriptionreminder.vercel.app/",
			"https://frontsubscriptionreminder.vercel.app",
		}),
		middleware.Trace(),
	); err != nil {
		log.Fatal("failed to run server", err)
	}
}
