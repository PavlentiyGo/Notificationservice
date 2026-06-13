package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	currencypb "github.com/PavlentiyGo/notification-service/proto/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService struct {
	currency map[string]float32
	rwMutex  sync.RWMutex
	currencypb.UnimplementedCurrencyServiceServer
}

func NewCurrencyService() *CurrencyService {

	currency := make(map[string]float32, 2)
	return &CurrencyService{
		currency: currency,
		rwMutex:  sync.RWMutex{},
	}
}

func (s *CurrencyService) GetCurrentCurrency(
	ctx context.Context,
	request *currencypb.GetCurrentCurrencyRequest,
) (*currencypb.GetCurrentCurrencyResponse, error) {
	s.rwMutex.RLock()

	usd, ok := s.currency["USD"]
	if !ok {
		return nil, status.Error(codes.Internal, "no usd currency")
	}
	eur, ok := s.currency["EUR"]
	if !ok {
		return nil, status.Error(codes.Internal, "no eur currency")
	}

	s.rwMutex.RUnlock()

	resp := &currencypb.GetCurrentCurrencyResponse{
		USD: usd,
		EUR: eur,
	}
	return resp, nil

}
func (s *CurrencyService) goWorker(ctx context.Context) {

	ticker := time.NewTicker(24 * time.Hour)
	err := s.parseCurrency()
	if err != nil {
		log.Println("failed to get currency ", err.Error())
	}
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err = s.parseCurrency()
				if err != nil {
					log.Println("failed to get currency ", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
func (s *CurrencyService) parseCurrency() error {
	client := http.Client{}
	resp, err := client.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	var decode CurrencyResponse

	if err = json.NewDecoder(resp.Body).Decode(&decode); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.currency["USD"] = decode.Valute["USD"].Value
	s.currency["EUR"] = decode.Valute["EUR"].Value
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":5052")
	if err != nil {
		log.Fatalln("failed to create listener")
	}
	defer listener.Close()

	grpcServ := grpc.NewServer()
	service := NewCurrencyService()

	currencypb.RegisterCurrencyServiceServer(grpcServ, service)

	service.goWorker(context.Background())

	if err = grpcServ.Serve(listener); err != nil {
		log.Fatalln("failed to start serve")
	}
}
