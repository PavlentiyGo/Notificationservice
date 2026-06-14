package handlers

import (
	"net/http"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/config"
	middleware "github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/middlewares"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/request"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/response"
	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/http/server"
	"google.golang.org/grpc"
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

func (h *AnalysisHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
			Middlewares: []middleware.Middleware{
				middleware.Authorize(h.config.BotToken),
			},
		},
	}
}
