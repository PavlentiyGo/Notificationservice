package service

import (
	"context"
	"fmt"

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

	if err := subscription.Validate(); err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to validate subscription: %w", err)
	}
	subscriptionCreated, err := s.repo.CreateSubscriptions(
		ctx,
		subscription,
	)
	if err != nil {
		return domain.Subscription{}, err
	}
	return subscriptionCreated, nil

}
func (s *SubscriptionService) GetSubscription(
	ctx context.Context,
	userId int32,
) ([]domain.Subscription, error) {

	subscriptions, err := s.repo.GetSubscriptions(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (s *SubscriptionService) PatchSubscription(
	ctx context.Context,
	patch domain.SubscriptionPatch,
) (domain.Subscription, error) {

	subscription, err := s.repo.GetSubscriptionById(ctx, patch.ID)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to get subs id: %w", err)
	}
	patchedSubscription, err := patch.Apply(subscription)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to apply patch on subsc: %w", err)
	}
	if err = s.repo.PatchSubscription(ctx, patchedSubscription); err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to patch subscription: %w", err)
	}
	return patchedSubscription, nil

}
