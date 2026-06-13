package handler

import (
	"context"
	"log"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	currencypb "github.com/PavlentiyGo/notification-service/proto/currency"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AnalysisHandler struct {
	currencyClient currencypb.CurrencyServiceClient
	service        *service.AnalysisService
	analysispb.UnimplementedAnalysisServiceServer
}

func NewAnalysisHandler(
	conn *grpc.ClientConn,
	service *service.AnalysisService,
) *AnalysisHandler {
	return &AnalysisHandler{
		currencyClient: currencypb.NewCurrencyServiceClient(conn),
		service:        service,
	}
}

func (h *AnalysisHandler) GetStatistics(
	ctx context.Context,
	req *analysispb.GetStatisticsRequest,
) (*analysispb.GetStatisticsResponse, error) {

	payments, err := h.service.GetStatistics(
		ctx,
		req.UserId,
		req.ThisMonth,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user payments: "+err.Error())
	}
	resp, err := h.currencyClient.GetCurrentCurrency(
		ctx,
		&currencypb.GetCurrentCurrencyRequest{},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get current currency to convert values: "+err.Error())
	}

	currency := domain.Currency{
		MainCurrency: req.TotalCurrency.String(),
		EUR:          float64(resp.EUR),
		USD:          float64(resp.USD),
	}

	groupedPayments := h.service.GroupPayments(
		ctx,
		payments,
		currency,
	)
	respProto := StatisticResponse(groupedPayments)
	log.Println(respProto)
	return respProto, nil
}
