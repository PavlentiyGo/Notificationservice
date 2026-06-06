package handlers

import (
	"fmt"
	"net/http"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/request"
	http2 "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/response"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubscriptionHandler struct {
	client subscriptionpb.SubscriptionServiceClient
}

func NewSubscriptionHandler(conn *grpc.ClientConn) *SubscriptionHandler {
	return &SubscriptionHandler{
		client: subscriptionpb.NewSubscriptionServiceClient(conn),
	}
}

func (h *SubscriptionHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method: http.MethodGet,
		},
	}

}

func (h *SubscriptionHandler) CreateSubscription(
	w http.ResponseWriter,
	r *http.Request,
) {
	responseHandler := http2.NewResponseHandler(w)

	user, err := middleware.UserFromCtx(r.Context())
	if err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusInternalServerError)
		return
	}

	var subscriptionRequest CreateSubscriptionRequest

	if err = request.DecodeAndValidate(r, &subscriptionRequest); err != nil {
		responseHandler.ErrorResponse(fmt.Sprintf("invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	currency, ok := subscriptionpb.Currency_value[subscriptionRequest.Currency]
	if !ok {
		responseHandler.ErrorResponse(fmt.Sprintf("invalid currency field: %s", subscriptionRequest.Currency), http.StatusBadRequest)
		return
	}
	subType, ok := subscriptionpb.SubscriptionType_value[subscriptionRequest.Type]
	if !ok {
		responseHandler.ErrorResponse(fmt.Sprintf("invalid subType field: %s", subscriptionRequest.Type), http.StatusBadRequest)
		return
	}
	resp, err := h.client.CreateSubscription(r.Context(), &subscriptionpb.CreateSubscriptionRequest{
		User: &subscriptionpb.User{
			Id:         user.ID,
			UserName:   user.Username,
			FirstName:  user.FirstName,
			SecondName: user.LastName,
		},
		Price:     subscriptionRequest.Price,
		Currency:  subscriptionpb.Currency(currency),
		Name:      subscriptionRequest.Name,
		Type:      subscriptionpb.SubscriptionType(subType),
		BillingAt: timestamppb.New(subscriptionRequest.BillingAt),
	})

}
