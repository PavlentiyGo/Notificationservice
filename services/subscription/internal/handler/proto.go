package handler

import (
	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func subscriptionToProto(
	val domain.Subscription,
) *subscriptionpb.CreateSubscriptionResponse {

	currency, _ := subscriptionpb.Currency_value[val.Currency]
	subType, _ := subscriptionpb.SubscriptionType_value[val.Type]
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

func subscriptionsToProto(
	values []domain.Subscription,
) []*subscriptionpb.CreateSubscriptionResponse {
	valuesProto := make([]*subscriptionpb.CreateSubscriptionResponse, 0, len(values))

	for _, val := range values {
		proto := subscriptionToProto(val)
		valuesProto = append(valuesProto, proto)
	}
	return valuesProto
}
