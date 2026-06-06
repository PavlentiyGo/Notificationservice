package service

import (
	"context"

	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(
	subRepo *repository.SubscriptionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		repo: subRepo,
	}
}

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	subscription domain.Subscription,
) (domain.Subscription, error) {

	subscriptionCreated, err := s.repo.CreateSubscriptions(
		ctx,
		subscription,
	)
	if err != nil {
		return domain.Subscription{}, err
	}
	return subscriptionCreated, nil

}
func (s *SubscriptionService) GetSubscription() {

}
