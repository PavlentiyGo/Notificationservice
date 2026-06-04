package main

import (
	"log"
	"net"

	pb "github.com/PavlentiyGo/notification-service/proto/notification"
	notification_handler "github.com/PavlentiyGo/notification-service/services/notification/handler"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	handler := notification_handler.NewNotificationHandler()

	pb.RegisterNotificationServiceServer(s, handler)

	log.Println("gRPC server listening on :50051")
	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
