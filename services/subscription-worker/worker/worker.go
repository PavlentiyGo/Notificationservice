package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	subscriptionpb "github.com/PavlentiyGo/notification-service/proto/subscription"
	"github.com/PavlentiyGo/notification-service/services/subscription-worker/publisher"
	"google.golang.org/grpc"
)

type Worker struct {
	publisher *publisher.Publisher
	client    subscriptionpb.SubscriptionServiceClient
}

func NewWorker(
	conn *grpc.ClientConn,
	publisher *publisher.Publisher,
) *Worker {
	return &Worker{
		client:    subscriptionpb.NewSubscriptionServiceClient(conn),
		publisher: publisher,
	}
}
func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour)

	subscriptions, err := w.GetSubscriptions(ctx)
	if err != nil {
		log.Printf("failed to get all subscriptions first time: %v", err)
		return err
	}
	w.CheckSubscriptions(ctx, subscriptions)
	go func() {
		defer ticker.Stop()

		for range ticker.C {
			subscriptions, err = w.GetSubscriptions(ctx)
			if err != nil {
				log.Printf("failed to get all subscription: %v", err)
				continue
			}
			w.CheckSubscriptions(ctx, subscriptions)
		}
	}()

	return nil
}

func (w *Worker) GetSubscriptions(ctx context.Context) ([]Subscription, error) {

	resp, err := w.client.GetAllSubscriptions(
		ctx,
		&subscriptionpb.GetAllSubscriptionsRequest{},
	)
	if err != nil {
		return nil, err
	}

	subscriptions := make([]Subscription, 0, len(resp.Subscriptions))

	for _, subscriptionGrpc := range resp.Subscriptions {
		subscription := Subscription{
			UserId:           subscriptionGrpc.User.Id,
			SubscriptionName: subscriptionGrpc.Name,
			BillingAt:        subscriptionGrpc.BillingAt.AsTime(),
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (w *Worker) CheckSubscriptions(
	ctx context.Context,
	subscriptions []Subscription,
) {

	for _, val := range subscriptions {
		timeDeadLine := time.Now().Add(48 * time.Hour)
		if timeDeadLine.After(val.BillingAt) && time.Now().Before(val.BillingAt) {
			w.publisher.PublishSubscriptionExpiring(
				ctx,
				val.UserId,
				val.SubscriptionName,
				getExpirationString(val.BillingAt),
			)
		}
	}
}

func getExpirationString(billingAt time.Time) string {

	duration := time.Now().Sub(billingAt)

	hours := int(duration.Hours())
	days := hours / 24

	var timeStr string
	if days > 0 {
		timeStr = fmt.Sprintf("%d дн.", days)
	} else {
		timeStr = fmt.Sprintf("%d ч.", hours)
	}

	return timeStr
}
