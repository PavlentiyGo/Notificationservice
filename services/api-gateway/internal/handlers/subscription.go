package handlers

import (
	"fmt"
	"net/http"
	"time"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/request"
	http2 "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/response"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubscriptionHandler struct {
	subscriptionClient subscriptionpb.SubscriptionServiceClient
	config             config.Config
}

func NewSubscriptionHandler(conn *grpc.ClientConn, config config.Config) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionClient: subscriptionpb.NewSubscriptionServiceClient(conn),
		config:             config,
	}
}

func (h *SubscriptionHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions",
			Handler: h.GetSubscriptions,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		}, {
			Method:  http.MethodPost,
			Path:    "/subscriptions",
			Handler: h.CreateSubscription,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		}, {
			Method:  http.MethodPatch, // TODO make
			Path:    "/subscriptions",
			Handler: nil,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		}, {
			Method:  http.MethodDelete, // TODO make
			Path:    "/subscriptions",
			Handler: nil,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
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
	billingAt, err := time.Parse("2006-01-02", subscriptionRequest.BillingAt)
	if err != nil {
		responseHandler.ErrorResponse(fmt.Sprintf("invalid billingAt field: %s", subscriptionRequest.BillingAt), http.StatusBadRequest)
		return
	}
	resp, err := h.subscriptionClient.CreateSubscription(r.Context(), &subscriptionpb.CreateSubscriptionRequest{
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
		BillingAt: timestamppb.New(billingAt),
	})
	if err != nil {
		responseHandler.GRPCErrorResponse(err)
		return
	}
	dtoResp := CreateSubscriptionResponse{
		SubscriptionId: resp.SubscriptionId,
		Price:          resp.Price,
		Currency:       resp.Currency.String(),
		Name:           resp.Name,
		Type:           resp.Type.String(),
		BillingAt:      resp.BillingAt.AsTime().Format("2006-01-02"),
	}
	responseHandler.JsonResponse(http.StatusCreated, dtoResp)

}

func (h *SubscriptionHandler) GetSubscriptions(
	w http.ResponseWriter,
	r *http.Request,
) {

	responseHandler := http2.NewResponseHandler(w)
	user, err := middleware.UserFromCtx(r.Context())
	if err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := h.subscriptionClient.GetSubscriptions(r.Context(), &subscriptionpb.GetSubscriptionsRequest{UserId: user.ID})
	if err != nil {
		responseHandler.GRPCErrorResponse(err)
		return
	}

	respDto := SubscriptionsDtoFromProto(resp)
	if len(respDto) == 0 {
		responseHandler.NoContentResponse()
		return
	}
	responseHandler.JsonResponse(http.StatusOK, respDto)
}
