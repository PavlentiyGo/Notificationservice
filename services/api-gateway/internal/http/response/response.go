package response

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseHandler struct {
	http.ResponseWriter
}

func NewResponseHandler(
	w http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		ResponseWriter: w,
	}
}

func (r *ResponseHandler) ErrorResponse(
	msg string,
	statusCode int,
) {
	r.JsonResponse(statusCode, map[string]string{
		"error": msg,
	})
}

func (r *ResponseHandler) JsonResponse(
	statusCode int,
	body any,
) {
	r.WriteHeader(statusCode)
	json.NewEncoder(r.ResponseWriter).Encode(&body)
}

func (r *ResponseHandler) GRPCErrorResponse(err error) {
	errMessage := status.Convert(err).Message()
	switch status.Code(err) {
	case codes.AlreadyExists:
		r.ErrorResponse(errMessage, http.StatusConflict)
	case codes.InvalidArgument:
		r.ErrorResponse(errMessage, http.StatusBadRequest)
	case codes.NotFound:
		r.ErrorResponse(errMessage, http.StatusNotFound)
	default:
		r.ErrorResponse("internal error", http.StatusInternalServerError)
	}
}
