package main

import (
	"context"
	"log"

	"github.com/PavlentiyGo/notification-service/services/subscription-worker/config"
	"github.com/PavlentiyGo/notification-service/services/subscription-worker/publisher"
	"github.com/PavlentiyGo/notification-service/services/subscription-worker/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg := config.NewConfigMust()

	publisher, err := publisher.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("failed to create new publisher: %v", err)
	}

	subscriptionConn, err := grpc.NewClient(cfg.SubscriptionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	worker := worker.NewWorker(subscriptionConn, publisher)

	if err = worker.Run(context.Background()); err != nil {
		log.Printf("failed to run worker properly: %v", err)
	}

}
