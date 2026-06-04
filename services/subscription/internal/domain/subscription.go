package domain

import "time"

type Subscription struct {
	SubscriptionId *int

	UserId    int
	Price     int
	Currency  string
	Name      string
	Type      string
	BillingAt time.Time
}

func (s *Subscription) Validate() {
	if s.UserId < 0 {
		
	}

}
