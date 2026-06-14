package repository

import (
	"context"
	"fmt"

	"github.com/PavlentiyGo/notification-service/services/analysis/internal/config"
	"github.com/PavlentiyGo/notification-service/services/analysis/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalysisRepository struct {
	pool   *pgxpool.Pool
	config config.Config
}

func NewAnalysisRepository(pool *pgxpool.Pool, cfg config.Config) *AnalysisRepository {
	return &AnalysisRepository{
		pool:   pool,
		config: cfg,
	}
}

func (r *AnalysisRepository) GetStatistics(
	ctx context.Context,
	userId int32,
	thisMonth bool,
) ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	SELECT 
    subscription_name, 
    subscription_type, 
    subscription_currency, 
    price
	FROM payments
	WHERE user_id = $1
	
	`
	if thisMonth {
		sqlQuery += `
		AND EXTRACT(MONTH FROM date) = EXTRACT(MONTH FROM NOW())
  		AND EXTRACT(YEAR FROM date) = EXTRACT(YEAR FROM NOW());
		`
	}

	rows, err := r.pool.Query(ctx, sqlQuery, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to query sql: %w", err)
	}
	payments := make([]PaymentModal, 0)

	for rows.Next() {
		var payment PaymentModal

		err = rows.Scan(
			&payment.SubscriptionName,
			&payment.SubscriptionType,
			&payment.SubscriptionCurrency,
			&payment.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		payments = append(payments, payment)
	}
	domainPayments := PaymentsDomainFromModals(payments)

	return domainPayments, err
}

func (r *AnalysisRepository) AddPayment(
	ctx context.Context,
	payment domain.Payment,
	userId int32,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.QueryTimeout)
	defer cancel()

	sqlQuery := `
	INSERT INTO payments(user_id,date, subscription_name,subscription_type, subscription_currency, price)
	VALUES($1,$2,$3,$4,$5,$6);
	`

	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		userId,
		payment.BillingAt,
		payment.SubscriptionName,
		payment.SubscriptionType,
		payment.SubscriptionCurrency,
		payment.Price,
	)
	if err != nil {
		return fmt.Errorf("failed to execute sql: %w", err)
	}

	return nil
}
