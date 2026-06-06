package main

import (
	"context"
	"log"
	"net"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/config"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/handler"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/repository"
	subscription_pool "github.com/PavlentiyGo/notification-service/services/subscription/internal/repository/pool"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/service"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.NewConfigMust()
	ctx := context.Background()

	pool, err := subscription_pool.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize subscription pool: %s", err)
	}
	subscriptionRepo := repository.NewSubscriptionRepository(pool, cfg)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	grpcServ := grpc.NewServer()
	subscriptionpb.RegisterSubscriptionServiceServer(grpcServ, subscriptionHandler)

	listener, err := net.Listen("tcp", ":"+cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen subscription serv: %s", err)
	}
	defer listener.Close()

	log.Printf("profile-service gRPC starting on %s", cfg.Addr)

	if err = grpcServ.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
