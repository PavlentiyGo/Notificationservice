package handlers

import "time"

type CreateSubscriptionRequest struct {
	Price     int32     `json:"price" validate:"required,numeric"`
	Currency  string    `json:"currency" validate:"oneof= RUB EUR USD"`
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	Type      string    `json:"type" validate:"oneof= STREAMING SOFTWARE UTILITIES FINANCE HEALTH EDUCATION OTHER"`
	BillingAt time.Time `json:"billing_at" validate:"required,datetime=2006-01-02"`
}
