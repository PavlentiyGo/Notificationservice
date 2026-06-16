package repository

import (
	"time"

	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
)

type Subscription struct {
	SubscriptionId int32

	UserId    int32
	Price     float64
	Currency  string
	Name      string
	Type      string
	BillingAt time.Time
}

func subscriptionDomainFromModal(subscriptions []Subscription) []domain.Subscription {

	subscriptionsDomain := make([]domain.Subscription, 0, len(subscriptions))

	for _, val := range subscriptions {
		subDomain := domain.Subscription{
			SubscriptionId: &val.SubscriptionId,
			User: domain.User{
				Id: val.UserId,
			},
			Price:     val.Price,
			Currency:  val.Currency,
			Name:      val.Name,
			Type:      val.Type,
			BillingAt: val.BillingAt,
		}
		subscriptionsDomain = append(subscriptionsDomain, subDomain)
	}

	return subscriptionsDomain
}
