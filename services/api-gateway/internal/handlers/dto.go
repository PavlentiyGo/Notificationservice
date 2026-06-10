package handlers

import (
	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
)

type CreateSubscriptionRequest struct {
	Price     int32  `json:"price" validate:"required,numeric"`
	Currency  string `json:"currency" validate:"oneof= RUB EUR USD"`
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Type      string `json:"type" validate:"oneof= STREAMING SOFTWARE UTILITIES FINANCE HEALTH EDUCATION OTHER"`
	BillingAt string `json:"billing_at" validate:"required,datetime=2006-01-02"`
}
type CreateSubscriptionResponse struct {
	SubscriptionId int32 `json:"subscription_id"`

	Price     int32  `json:"price"`
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	BillingAt string `json:"billing_at"`
}
type GetSubscriptionsResponse struct {
	Subscriptions []CreateSubscriptionResponse `json:"subscriptions"`
}

func SubscriptionsDtoFromProto(
	subscriptions *subscriptionpb.GetSubscriptionsResponse,
) []CreateSubscriptionResponse {
	subscDto := make([]CreateSubscriptionResponse, 0, len(subscriptions.Subscriptions))

	for _, val := range subscriptions.Subscriptions {
		dto := CreateSubscriptionResponse{
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
