package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PavlentiyGo/notification-service/services/subscription/internal/config"
	"github.com/PavlentiyGo/notification-service/services/subscription/internal/domain"
	subscription_errors "github.com/PavlentiyGo/notification-service/services/subscription/internal/errors"
	"github.com/jackc/pgx/v5"
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
	INSERT INTO users(id,user_name,first_name,second_name)  
	VALUES($1,$2,$3,$4);
	`

	_, err := r.Exec(
		ctx,
		sqlQuery,
		user.Id,
		user.UserName,
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

	if err := r.CreateUser(ctx, subscription.User); err != nil {
		return domain.Subscription{}, fmt.Errorf("failed to create user: %w", err)
	}
	sqlQuery := `
	INSERT INTO subscriptions(user_id,price,currency,name,type,billing_at)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING id;
	`

	row := r.QueryRow(
		ctx,
		sqlQuery,
		subscription.User.Id,
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
	userId int64,
) ([]domain.Subscription, error) {

	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	SELECT id,user_id,price,currency,name,type,billing_at
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
func (r *SubscriptionRepository) GetSubscriptionById(
	ctx context.Context,
	id int32,
) (domain.Subscription, error) {

	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	SELECT id,user_id,price,currency,name,type,billing_at
	FROM subscriptions
	WHERE id = $1;
	`

	rows := r.QueryRow(
		ctx,
		sqlQuery,
		id,
	)

	var subscriptionModal Subscription

	err := rows.Scan(
		&subscriptionModal.SubscriptionId,
		&subscriptionModal.UserId,
		&subscriptionModal.Price,
		&subscriptionModal.Currency,
		&subscriptionModal.Name,
		&subscriptionModal.Type,
		&subscriptionModal.BillingAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, fmt.Errorf("no such subscription: %w", subscription_errors.ErrNotFound)
		}
		return domain.Subscription{}, fmt.Errorf("failed to exec sql query: %w", err)
	}

	subscriptionsDomain := subscriptionDomainFromModal([]Subscription{subscriptionModal})

	return subscriptionsDomain[0], nil
}

func (r *SubscriptionRepository) PatchSubscription(
	ctx context.Context,
	subscription domain.Subscription,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	UPDATE subscriptions
	SET 
	price = $1,
	currency = $2,
	name = $3,
	type = $4,
	billing_at = $5
	WHERE id = $6;
	`

	_, err := r.Exec(
		ctx,
		sqlQuery,
		subscription.Price,
		subscription.Currency,
		subscription.Name,
		subscription.Type,
		subscription.BillingAt,
		subscription.SubscriptionId,
	)
	if err != nil {
		return fmt.Errorf("failed to exec sql query: %w", err)
	}

	return nil

}
