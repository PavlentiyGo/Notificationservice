package handler

import (
	"context"
	"errors"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	errors2 "github.com/PavlentiyGo/notification-service/services/subscription/internal/errors"
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
		User: domain.User{
			Id:         req.User.Id,
			UserName:   req.User.UserName,
			Name:       req.User.FirstName,
			SecondName: req.User.SecondName,
		},
		Price:     req.Price,
		Currency:  req.Currency.String(),
		Name:      req.Name,
		Type:      req.Type.String(),
		BillingAt: req.BillingAt.AsTime(),
	}
	createdSubscription, err := h.service.CreateSubscription(ctx, subscription)
	if err != nil {
		if errors.Is(err, errors2.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	subProto := subscriptionToProto(createdSubscription)

	return subProto, nil
}

func (h *SubscriptionHandler) GetSubscriptions(
	ctx context.Context,
	req *subscriptionpb.GetSubscriptionsRequest,
) (*subscriptionpb.GetSubscriptionsResponse, error) {

	subscriptions, err := h.service.GetSubscription(ctx, req.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoResp := subscriptionsToProto(subscriptions)

	resp := &subscriptionpb.GetSubscriptionsResponse{
		Subscriptions: protoResp,
	}

	return resp, nil
}
