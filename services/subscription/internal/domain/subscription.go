package domain

import (
	"fmt"
	"time"

	errors2 "github.com/PavlentiyGo/notification-service/services/subscription/internal/errors"
)

type Subscription struct {
	SubscriptionId *int32

	User      User
	Price     int32
	Currency  string
	Name      string
	Type      string
	BillingAt time.Time
}

func (s *Subscription) Validate() error {
	if s.Price <= 0 {
		return fmt.Errorf("wrong price for subscription: %d: %w", s.Price, errors2.ErrInvalidArgument)
	}
	if s.BillingAt.Before(time.Now()) {
		return fmt.Errorf("wrong billingAt for subscription: %s: %w", s.BillingAt.String(), errors2.ErrInvalidArgument)
	}
	return nil
}
