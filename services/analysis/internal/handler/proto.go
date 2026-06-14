package handler

import (
	"fmt"

	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
)

func StatisticResponse(payments map[string]domain.GroupPayment) *analysispb.GetStatisticsResponse {

	resp := &analysispb.GetStatisticsResponse{
		TotalSum:     "",
		PaymentsList: []*analysispb.PaymentList{},
	}
	var totalSum float64
	for key, val := range payments {
		paymentType := analysispb.SubscriptionType_value[key]
		subscriptionNames := make([]string, 0, len(val.Payments))

		paymentList := &analysispb.PaymentList{
			PaymentsType:      analysispb.SubscriptionType(paymentType),
			TotalSum:          fmt.Sprintf("%.2f", val.TotalPrice),
			SubscriptionsName: subscriptionNames,
		}
		totalSum += val.TotalPrice
		for _, payment := range val.Payments {
			paymentList.SubscriptionsName = append(paymentList.SubscriptionsName, payment.SubscriptionName)
		}
		resp.PaymentsList = append(resp.PaymentsList, paymentList)
	}
	resp.TotalSum = fmt.Sprintf("%.2f", totalSum)
	return resp
}
