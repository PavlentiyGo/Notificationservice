package repository

import (
	"context"
	"fmt"

	"github.com/PavlentiyGo/notification-service/services/subscription/internal/config"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	*pgxpool.Pool
	config config.Config
}

func NewSubscriptionRepository(
	pool *pgxpool.Pool,
	config config.Config,
) *SubscriptionRepository {
	return &SubscriptionRepository{
		Pool:   pool,
		config: config,
	}
}

func (r *SubscriptionRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) error {

	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	INSERT INTO users(user_id,name,second_name) 
	VALUES($1,$2,$3);
	`

	_, err := r.Exec(
		ctx,
		sqlQuery,
		user.Id,
		user.Name,
		user.SecondName,
	)
	if err != nil {
		if IsPgErr(err, "23505") {
			return nil
		}
		return fmt.Errorf("failed to exec sql query: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) CreateSubscriptions(
	ctx context.Context,
	subscription domain.Subscription,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	INSERT INTO subscriptions(user_id,price,currency,name,type,billing_at)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING subscription_id;
	`

	row := r.QueryRow(
		ctx,
		sqlQuery,
		subscription.UserId,
		subscription.Price,
		subscription.Currency,
		subscription.Name,
		subscription.Type,
		subscription.BillingAt,
	)

	var sub Subscription
	if err := row.Scan(
		&sub.SubscriptionId,
	); err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}
	subscription.SubscriptionId = &sub.SubscriptionId
	return subscription, nil
}

func (r *SubscriptionRepository) GetSubscriptions(
	ctx context.Context,
	userId int,
) ([]domain.Subscription, error) {

	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	SELECT subscription_id,user_id,price,currency,name,type,billing_at
	FROM subscriptions
	WHERE user_id = $1;	
	`

	rows, err := r.Query(
		ctx,
		sqlQuery,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exec sql query: %w", err)
	}

	var subscriptionsModal []Subscription

	for rows.Next() {
		var subscription Subscription
		err = rows.Scan(
			&subscription.SubscriptionId,
			&subscription.UserId,
			&subscription.Price,
			&subscription.Currency,
			&subscription.Name,
			&subscription.Type,
			&subscription.BillingAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		subscriptionsModal = append(subscriptionsModal, subscription)
	}
	subscriptionsDomain := subscriptionDomainFromModal(subscriptionsModal)

	return subscriptionsDomain, nil
}
