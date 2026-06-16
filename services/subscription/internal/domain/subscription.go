package domain

import (
	"fmt"
	"slices"
	"time"

	errors2 "github.com/PavlentiyGo/notification-service/services/subscription/internal/errors"
)

var Currencies = []string{
	"RUB", "USD", "EUR",
}
var Types = []string{
	"STREAMING", // Нетфликс, Яндекс Плюс, Спотифай
	"SOFTWARE",  // JetBrains, ChatGPT Plus, Облако
	"UTILITIES", // ЖКХ, интернет, мобильная связь
	"FINANCE",   // Кредиты, страховки, банковское обслуживание
	"HEALTH",    // Фитнес-клуб, доставка здоровой еды, витамины
	"EDUCATION", // Курсы, онлайн-школы, подписки на книги
	"OTHER",
}

type Subscription struct {
	SubscriptionId *int32

	User      User
	Price     float64
	Currency  string
	Name      string
	Type      string
	BillingAt time.Time
}

func (s *Subscription) Validate() error {
	if s.Price <= 0 {
		return fmt.Errorf("wrong price for subscription: %f: %w", s.Price, errors2.ErrInvalidArgument)
	}
	if index := slices.Index(Currencies, s.Currency); index == -1 {
		return fmt.Errorf("wrong currency for subscription: %s: %w", s.Currency, errors2.ErrInvalidArgument)
	}
	if index := slices.Index(Types, s.Type); index == -1 {
		return fmt.Errorf("wrong type for subscription: %s: %w", s.Type, errors2.ErrInvalidArgument)
	}
	if len(s.Name) < 3 || len(s.Name) > 100 {
		return fmt.Errorf("wrong name len for subscription: %s: %w", s.Name, errors2.ErrInvalidArgument)
	}
	if s.BillingAt.Before(time.Now()) {
		return fmt.Errorf("wrong billingAt for subscription: %s: %w", s.BillingAt.String(), errors2.ErrInvalidArgument)
	}
	return nil
}

type SubscriptionPatch struct {
	ID int32

	Price     *float64
	Currency  *string
	Name      *string
	Type      *string
	BillingAt *time.Time
}

func (p *SubscriptionPatch) Apply(
	subscription Subscription,
) (Subscription, error) {

	if p.Name != nil {
		subscription.Name = *p.Name
	}
	if p.Price != nil {
		subscription.Price = *p.Price
	}
	if p.Currency != nil {
		subscription.Currency = *p.Currency
	}
	if p.Type != nil {
		subscription.Type = *p.Type
	}
	if p.BillingAt != nil {
		subscription.BillingAt = *p.BillingAt
	}
	if err := subscription.Validate(); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}
