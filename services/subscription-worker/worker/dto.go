package worker

import "time"

type Subscription struct {
	UserId int64

	SubscriptionName string
	BillingAt        time.Time
}
