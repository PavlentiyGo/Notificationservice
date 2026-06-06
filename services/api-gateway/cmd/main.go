package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	grpc_client "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/grpc"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/handlers"
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

	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionClient)
	serv.ChainRoutes(subscriptionHandler.Routes()...)

}
