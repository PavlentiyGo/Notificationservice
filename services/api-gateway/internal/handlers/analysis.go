package handlers

import (
	"net/http"
	"time"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	middleware "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/request"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/response"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalysisHandler struct {
	config         config.Config
	analysisClient analysispb.AnalysisServiceClient
}

func NewAnalysisHandler(
	cfg config.Config,
	conn *grpc.ClientConn,
) *AnalysisHandler {
	return &AnalysisHandler{
		config:         cfg,
		analysisClient: analysispb.NewAnalysisServiceClient(conn),
	}
}

func (h *AnalysisHandler) GetStatistics(
	w http.ResponseWriter,
	r *http.Request,
) {
	responseHandler := response.NewResponseHandler(w)

	user, err := middleware.UserFromCtx(r.Context())
	if err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusInternalServerError)
		return
	}
	var requestHTTP GetStatisticRequest
	if err = request.DecodeAndValidate(r, &requestHTTP); err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusBadRequest)
		return
	}

	totalCurrency := analysispb.Currency_value[requestHTTP.TotalCurrency]
	respGRPC, err := h.analysisClient.GetStatistics(
		r.Context(),
		&analysispb.GetStatisticsRequest{
			TotalCurrency: analysispb.Currency(totalCurrency),
			UserId:        user.ID,
			ThisMonth:     requestHTTP.ThisMonth,
		},
	)
	if err != nil {
		responseHandler.GRPCErrorResponse(err)
		return
	}
	responseHttp := StatisticResponseFromProto(respGRPC)

	responseHandler.JsonResponse(http.StatusOK, responseHttp)
}

func (h *AnalysisHandler) AddPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	responseHandler := response.NewResponseHandler(w)

	user, err := middleware.UserFromCtx(r.Context())
	if err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusInternalServerError)
		return
	}
	var requestHTTP AddPaymentRequest
	if err = request.DecodeAndValidate(r, &requestHTTP); err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusBadRequest)
		return
	}

	subType := analysispb.SubscriptionType_value[requestHTTP.Type]
	subCurrency := analysispb.Currency_value[requestHTTP.Currency]
	billingAtTime, err := time.Parse("2006-01-2", requestHTTP.BillingAt)
	if err != nil {
		responseHandler.ErrorResponse(err.Error(), http.StatusBadRequest)
		return
	}
	addPaymentRequest := &analysispb.AddPaymentRequest{
		UserId:         user.ID,
		SubscriptionId: requestHTTP.SubscriptionId,
		BillingAt:      timestamppb.New(billingAtTime),
		Type:           analysispb.SubscriptionType(subType),
		Name:           requestHTTP.Name,
		Currency:       analysispb.Currency(subCurrency),
		Price:          requestHTTP.Price,
	}
	respGRPC, err := h.analysisClient.AddPayment(
		r.Context(),
		addPaymentRequest,
	)
	if err != nil {
		responseHandler.GRPCErrorResponse(err)
		return
	}
	respHTTP := struct {
		BillingAt string `json:"next_billing_at"`
	}{
		BillingAt: respGRPC.NextBillingAt.AsTime().Format("2006-01-02"),
	}
	responseHandler.JsonResponse(http.StatusOK, respHTTP)

}

// 998998854

func (h *AnalysisHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/statistics",
			Handler: h.GetStatistics,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		}, {
			Method:  http.MethodPost,
			Path:    "/payments",
			Handler: h.AddPayment,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		},
	}
}
