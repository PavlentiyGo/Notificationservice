package handler

import (
	"context"
	"log"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionHandler struct {
	subscriptionpb.UnimplementedSubscriptionServiceServer
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) CreateSubscription(
	ctx context.Context,
	req *subscriptionpb.CreateSubscriptionRequest,
) (*subscriptionpb.CreateSubscriptionResponse, error) {

	subscription := domain.Subscription{
		SubscriptionId: nil,
		UserId:         req.User.Id,
		Price:          req.Price,
		Currency:       req.Currency.String(),
		Name:           req.Name,
		Type:           req.Type.String(),
		BillingAt:      req.BillingAt.AsTime(),
	}
	createdSubscription, err := h.service.CreateSubscription(ctx, subscription)
	if err != nil {
		log.Print(err, subscription.Currency)
		return nil, status.Error(codes.Internal, err.Error())
	}
	subProto := subscriptionsToProto(createdSubscription)

	return subProto, nil
}

func (h *SubscriptionHandler) GetSubscriptions(
	ctx context.Context,
	req *subscriptionpb.GetSubscriptionsRequest,
) (*subscriptionpb.GetSubscriptionsResponse, error) {
	return nil, nil
}
