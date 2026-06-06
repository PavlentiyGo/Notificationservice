package handler

import (
	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func subscriptionsToProto(
	val domain.Subscription,
) *subscriptionpb.CreateSubscriptionResponse {

	currency, _ := subscriptionpb.Currency_value[val.Currency]
	subType, _ := subscriptionpb.Currency_value[val.Type]
	billingAt := timestamppb.New(val.BillingAt)

	subProto := &subscriptionpb.CreateSubscriptionResponse{
		SubscriptionId: *val.SubscriptionId,

		Price:     val.Price,
		Currency:  subscriptionpb.Currency(currency),
		Name:      val.Name,
		Type:      subscriptionpb.SubscriptionType(subType),
		BillingAt: billingAt,
	}

	return subProto
}
