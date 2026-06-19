package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/PavlentiyGo/notification-service/services/notification/config"
	"github.com/PavlentiyGo/notification-service/services/notification/consumer"
	"github.com/PavlentiyGo/notification-service/services/notification/sender"
)

func main() {

	cfg := config.NewConfigMust()

	Sender := sender.NewSender(cfg)

	Consumer, err := consumer.NewConsumer(cfg, Sender)
	if err != nil {
		log.Fatal("failed to create consumer" + err.Error())
	}
	defer Consumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err = Consumer.Run(ctx); err != nil {
		log.Fatalf("consumer error: %v", err)
	}

}
