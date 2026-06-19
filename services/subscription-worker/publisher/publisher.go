package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PavlentiyGo/notification-service/services/subscription-worker/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(cfg config.Config) (*Publisher, error) {

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbitmq conn: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbitmq channel: %w", err)
	}
	for _, queue := range []string{"subscription.expiring"} {
		_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create queue for rabbitmq: %s: %w", queue, err)
		}
	}

	return &Publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Publisher) Close() {
	p.conn.Close()
	p.channel.Close()
}
func (p *Publisher) PublishSubscriptionExpiring(
	ctx context.Context,
	userId int64,
	subscriptionName string,
	expiringAt string,
) {
	p.publish(ctx,
		"subscription.expiring",
		map[string]any{
			"userId":           userId,
			"subscriptionName": subscriptionName,
			"expiringAt":       expiringAt,
		},
	)
}

func (p *Publisher) publish(
	ctx context.Context,
	queue string,
	payload map[string]any,
) {
	event := Event{
		Type:    queue,
		Payload: payload,
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}
	err = p.channel.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         data,
		})
	if err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
