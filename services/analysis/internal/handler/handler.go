package handler

import (
	"context"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	currencypb "github.com/PavlentiyGo/notification-service/proto/currency"
	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalysisHandler struct {
	currencyClient     currencypb.CurrencyServiceClient
	subscriptionClient subscriptionpb.SubscriptionServiceClient
	service            *service.AnalysisService
	analysispb.UnimplementedAnalysisServiceServer
}

func NewAnalysisHandler(
	conn *grpc.ClientConn,
	conn2 *grpc.ClientConn,
	service *service.AnalysisService,
) *AnalysisHandler {
	return &AnalysisHandler{
		currencyClient:     currencypb.NewCurrencyServiceClient(conn),
		subscriptionClient: subscriptionpb.NewSubscriptionServiceClient(conn2),
		service:            service,
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
	return respProto, nil
}
func (h *AnalysisHandler) AddPayment(
	ctx context.Context,
	request *analysispb.AddPaymentRequest,
) (*analysispb.AddPaymentResponse, error) {

	billingAtTime := request.BillingAt.AsTime()
	nextBillingAt, err := h.service.AddPayment(
		ctx,
		domain.Payment{
			BillingAt:            &billingAtTime,
			SubscriptionName:     request.Name,
			SubscriptionType:     request.Type.String(),
			SubscriptionCurrency: request.Currency.String(),
			Price:                request.Price,
		},
		request.UserId,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	_, err = h.subscriptionClient.PatchSubscription(
		ctx,
		&subscriptionpb.PatchSubscriptionRequest{
			SubscriptionId: request.SubscriptionId,
			BillingAt:      timestamppb.New(nextBillingAt),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to perform patch"+err.Error())
	}

	resp := &analysispb.AddPaymentResponse{
		NextBillingAt: timestamppb.New(nextBillingAt),
	}

	return resp, nil
}
