package domain

import "time"

type Payment struct {
	BillingAt *time.Time

	SubscriptionName     string
	SubscriptionType     string
	SubscriptionCurrency string
	Price                float64
}

type GroupPayment struct {
	Payments   []Payment
	TotalPrice float64
}
