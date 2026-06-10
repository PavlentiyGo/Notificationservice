package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type CurrencyService struct {
}

func main() {
	listener, err := net.Listen("tcp", "5052")
	if err != nil {
		log.Fatalln("failed to create listener")
	}
	defer listener.Close()

	grpcServ := grpc.NewServer()

	grpcServ.Serve(listener)

}
