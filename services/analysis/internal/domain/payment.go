package domain

type Payment struct {
	SubscriptionName     string
	SubscriptionType     string
	SubscriptionCurrency string
	Price                float64
}

type GroupPayment struct {
	Payments   []Payment
	TotalPrice float64
}
