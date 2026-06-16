package handler

import (
	"context"
	"errors"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	subscription_errors "github.com/PavlentiyGo/notification-service/services/subscription/internal/errors"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		if errors.Is(err, subscription_errors.ErrInvalidArgument) {
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

func (h *SubscriptionHandler) PatchSubscription(
	ctx context.Context,
	req *subscriptionpb.PatchSubscriptionRequest,
) (*subscriptionpb.PathSubscriptionResponse, error) {

	subscriptionPatch := domain.SubscriptionPatch{
		ID:    req.SubscriptionId,
		Price: req.Price,
		Name:  req.Name,
	}
	if req.Currency != nil {
		currency := req.Currency.String()
		subscriptionPatch.Currency = &currency
	}
	if req.Type != nil {
		subType := req.Type.String()
		subscriptionPatch.Type = &subType
	}
	if req.BillingAt != nil {
		billingAt := req.BillingAt.AsTime()
		subscriptionPatch.BillingAt = &billingAt
	}
	patchedSubscription, err := h.service.PatchSubscription(ctx, subscriptionPatch)
	if err != nil {
		if errors.Is(err, subscription_errors.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, subscription_errors.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	currency := subscriptionpb.Currency_value[patchedSubscription.Currency]
	subType := subscriptionpb.SubscriptionType_value[patchedSubscription.Type]
	return &subscriptionpb.PathSubscriptionResponse{
		SubscriptionId: *patchedSubscription.SubscriptionId,
		Price:          patchedSubscription.Price,
		Currency:       subscriptionpb.Currency(currency),
		Name:           patchedSubscription.Name,
		Type:           subscriptionpb.SubscriptionType(subType),
		BillingAt:      timestamppb.New(patchedSubscription.BillingAt),
	}, nil
}
