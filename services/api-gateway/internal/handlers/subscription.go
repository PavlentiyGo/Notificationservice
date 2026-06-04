package handlers

import (
	"net/http"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"google.golang.org/grpc"
)

type SubscriptionHandler struct {
	client subscriptionpb.SubscriptionServiceClient
}

func NewSubscriptionHandler(conn *grpc.ClientConn) *SubscriptionHandler {
	return &SubscriptionHandler{
		client: subscriptionpb.NewSubscriptionServiceClient(conn),
	}
}

func (h *SubscriptionHandler) CreateSubscription(
	w http.ResponseWriter,
	r *http.Request,
) {

	h.client.


}
