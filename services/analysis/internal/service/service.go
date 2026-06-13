package service

import (
	"context"
	"math"

	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/repository"
)

type AnalysisService struct {
	repository *repository.AnalysisRepository
}

func NewAnalysisService(repo *repository.AnalysisRepository) *AnalysisService {
	return &AnalysisService{
		repository: repo,
	}
}

func (s *AnalysisService) GetStatistics(
	ctx context.Context,
	userId int32,
	thisMonth bool,
) ([]domain.Payment, error) {
	return s.repository.GetStatistics(ctx, userId, thisMonth)
}

func (s *AnalysisService) GroupPayments(
	ctx context.Context,
	payments []domain.Payment,
	currentCurrency domain.Currency,
) map[string]domain.GroupPayment {

	groupedPayments := make(map[string]domain.GroupPayment)

	for _, payment := range payments {
		val, ok := groupedPayments[payment.SubscriptionType]
		if !ok {
			groupedPayments[payment.SubscriptionType] = domain.GroupPayment{
				Payments: []domain.Payment{payment},
				TotalPrice: convertCurrency(
					payment.SubscriptionCurrency,
					currentCurrency.MainCurrency,
					currentCurrency,
					payment.Price,
				),
			}
		} else {
			val.Payments = append(val.Payments, payment)
			val.TotalPrice += convertCurrency(
				payment.SubscriptionCurrency,
				currentCurrency.MainCurrency,
				currentCurrency,
				payment.Price,
			)
			groupedPayments[payment.SubscriptionType] = val
		}
	}
	return groupedPayments
}
func convertCurrency(from, to string, currency domain.Currency, price float64) float64 {
	if from == to {
		return price
	}
	switch to {
	case "RUB":
		if from == "USD" {
			return math.Round(price*currency.USD*100) / 100
		} else if from == "EUR" {
			return math.Round(price*currency.EUR*100) / 100
		}
	case "USD":
		if from == "RUB" {
			return math.Round(price/currency.USD*100) / 100
		} else if from == "EUR" {
			return math.Round(price*currency.EUR/currency.USD*100) / 100
		}
	case "EUR":
		if from == "USD" {
			return math.Round(price*currency.USD/currency.EUR*100) / 100
		} else if from == "RUB" {
			return math.Round(price/currency.EUR*100) / 100
		}
	}
	return price
}
