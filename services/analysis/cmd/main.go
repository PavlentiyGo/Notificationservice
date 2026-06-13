package main

import (
	"context"
	"log"
	"net"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/config"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/handler"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/repository"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/repository/pool"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg := config.NewConfigMust()

	pool, err := pool.NewPool(context.Background(), cfg)
	if err != nil {
		log.Fatalln("failed to start pool " + err.Error())
	}
	analysisRepository := repository.NewAnalysisRepository(pool, cfg)
	analysisService := service.NewAnalysisService(analysisRepository)

	currencyConn, err := grpc.NewClient(cfg.CurrencyAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("failed to create currency client ", err.Error())
		return
	}
	defer currencyConn.Close()

	analysisHandler := handler.NewAnalysisHandler(currencyConn, analysisService)

	serv := grpc.NewServer()
	analysispb.RegisterAnalysisServiceServer(serv, analysisHandler)

	listener, err := net.Listen("tcp", cfg.AnalysisAddr)
	if err != nil {
		log.Println("failed to create tcp listener ", err.Error())
		return
	}
	if err = serv.Serve(listener); err != nil {
		log.Println("failed to serve grpc server: ", err.Error())
	}
}
