package handler

import (
	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
)

func StatisticResponse(payments map[string]domain.GroupPayment) *analysispb.GetStatisticsResponse {

	resp := &analysispb.GetStatisticsResponse{
		TotalSum:     0,
		PaymentsList: []*analysispb.PaymentList{},
	}

	for key, val := range payments {
		paymentType := analysispb.SubscriptionType_value[key]
		subscriptionNames := make([]string, len(val.Payments))

		paymentList := &analysispb.PaymentList{
			PaymentsType:      analysispb.SubscriptionType(paymentType),
			TotalSum:          0,
			SubscriptionsName: subscriptionNames,
		}
		for _, payment := range val.Payments {
			resp.TotalSum += payment.Price
			paymentList.TotalSum += payment.Price
			paymentList.SubscriptionsName = append(paymentList.SubscriptionsName, payment.SubscriptionName)
		}
		resp.PaymentsList = append(resp.PaymentsList, paymentList)
	}
	return resp
}
