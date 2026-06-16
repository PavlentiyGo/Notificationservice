package handlers

import (
	analysispb "github.com/PavlentiyGo/notification-service/proto/analysis"
	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
)

type CreateSubscriptionRequest struct {
	Price     float64 `json:"price" validate:"required,numeric"`
	Currency  string  `json:"currency" validate:"oneof= RUB EUR USD"`
	Name      string  `json:"name" validate:"required,min=3,max=100"`
	Type      string  `json:"type" validate:"oneof= STREAMING SOFTWARE UTILITIES FINANCE HEALTH EDUCATION OTHER"`
	BillingAt string  `json:"billing_at" validate:"required,datetime=2006-01-02"`
}
type SubscriptionResponse struct {
	SubscriptionId int32 `json:"subscription_id"`

	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	BillingAt string  `json:"billing_at"`
}
type GetSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
}

func SubscriptionsDtoFromProto(
	subscriptions *subscriptionpb.GetSubscriptionsResponse,
) []SubscriptionResponse {
	subscDto := make([]SubscriptionResponse, 0, len(subscriptions.Subscriptions))

	for _, val := range subscriptions.Subscriptions {
		dto := SubscriptionResponse{
			SubscriptionId: val.SubscriptionId,

			Price:     val.Price,
			Currency:  val.Currency.String(),
			Name:      val.Name,
			Type:      val.Type.String(),
			BillingAt: val.BillingAt.AsTime().Format("2006-01-02"),
		}
		subscDto = append(subscDto, dto)
	}
	return subscDto
}

type GetStatisticRequest struct {
	TotalCurrency string `json:"total_currency" validate:"oneof= RUB EUR USD"`
	ThisMonth     bool   `json:"this_month"`
}

type GetStatisticResponse struct {
	TotalSum string         `json:"total_sum"`
	Payments []PaymentTypes `json:"payments"`
}

type PaymentTypes struct {
	Type          string   `json:"type"`
	TotalSum      string   `json:"total_sum"`
	Subscriptions []string `json:"subscriptions"`
}

func StatisticResponseFromProto(
	response *analysispb.GetStatisticsResponse,
) GetStatisticResponse {

	payments := make([]PaymentTypes, 0, len(response.PaymentsList))
	statisticResponse := GetStatisticResponse{
		TotalSum: response.TotalSum,
		Payments: nil,
	}

	for _, payment := range response.PaymentsList {
		paymentType := PaymentTypes{
			Type:          payment.PaymentsType.String(),
			TotalSum:      payment.TotalSum,
			Subscriptions: payment.SubscriptionsName,
		}
		payments = append(payments, paymentType)
	}
	statisticResponse.Payments = payments
	return statisticResponse
}

type AddPaymentRequest struct {
	Currency string  `json:"currency" validate:"oneof= RUB EUR USD"`
	Date     *string `json:"date" validate:"omitempty,datetime=2006-01-02"`
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Price    float64 `json:"price" validate:"required,numeric"`
	Type     string  `json:"type" validate:"oneof= STREAMING SOFTWARE UTILITIES FINANCE HEALTH EDUCATION OTHER"`
}

type PatchSubscriptionRequest struct {
	Id int32 `json:"id"`

	Price     *float64 `json:"price"      validate:"omitempty,numeric"`
	Currency  *string  `json:"currency"   validate:"omitempty,oneof= RUB EUR USD"`
	Name      *string  `json:"name"       validate:"omitempty,min=3,max=100"`
	Type      *string  `json:"type"       validate:"omitempty,oneof= STREAMING SOFTWARE UTILITIES FINANCE HEALTH EDUCATION OTHER"`
	BillingAt *string  `json:"billing_at" validate:"omitempty,datetime=2006-01-02"`
}
