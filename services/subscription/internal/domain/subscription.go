package domain

import "time"

type Subscription struct {
	SubscriptionId *int32

	UserId    int32
	Price     int32
	Currency  string
	Name      string
	Type      string
	BillingAt time.Time
}

func (s *Subscription) Validate() {
	if s.UserId < 0 {

	}

}
