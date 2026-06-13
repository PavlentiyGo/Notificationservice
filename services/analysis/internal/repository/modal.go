package repository

import "github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"

type PaymentModal struct {
	SubscriptionName     string
	SubscriptionType     string
	SubscriptionCurrency string
	Price                float64
}

func PaymentsDomainFromModals(
	payments []PaymentModal,
) []domain.Payment {

	paymentsDomain := make([]domain.Payment, 0, len(payments))

	for _, val := range payments {
		paymentDomain := domain.Payment{
			SubscriptionName:     val.SubscriptionName,
			SubscriptionType:     val.SubscriptionType,
			SubscriptionCurrency: val.SubscriptionCurrency,
			Price:                val.Price,
		}
		paymentsDomain = append(paymentsDomain, paymentDomain)
	}

	return paymentsDomain
}
