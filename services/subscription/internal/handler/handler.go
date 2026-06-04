package handler

import (
	"context"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
)

type SubscriptionHandler struct {
	subscriptionpb.UnimplementedSubscriptionServiceServer
}

func (h *SubscriptionHandler) CreateSubscription(
	ctx context.Context,
	req *subscriptionpb.CreateSubscriptionRequest,
) (*subscriptionpb.CreateSubscriptionResponse, error) {

}

func (h *SubscriptionHandler) GetSubscriptions(
	ctx context.Context,
	req *subscriptionpb.GetSubscriptionsRequest,
) (*subscriptionpb.GetSubscriptionsResponse, error) {

}
