package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/PavlentiyGo/notification-service/services/notification/config"
	"github.com/PavlentiyGo/notification-service/services/notification/sender"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type    string
	Payload map[string]any
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	sender  *sender.Sender
}

func NewConsumer(
	cfg config.Config,
	sender *sender.Sender,
) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	for _, queue := range []string{"subscription.expiring"} {
		_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue: %s: %w", queue, err)
		}
	}
	return &Consumer{
		conn:    conn,
		channel: ch,
		sender:  sender,
	}, nil
}

func (c *Consumer) Close() {
	c.conn.Close()
	c.channel.Close()
}

func (c *Consumer) Run(ctx context.Context) error {
	expiringMsgs, err := c.channel.Consume("subscription.expiring", "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume subscrption.expiring channel: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-expiringMsgs:
			if !ok {
				return nil
			}
			c.handleSubscriptionExpiring(msg)
		}
	}
}

func (c *Consumer) handleSubscriptionExpiring(
	msg amqp.Delivery,
) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("failed to unmarshal event for subscription.expiring %v", err)
		msg.Nack(false, false)
		return
	}

	userId, _ := event.Payload["userId"].(int64)
	subscriptionName, _ := event.Payload["subscriptionName"].(string)
	expiresIn, _ := event.Payload["expiresIn"].(string)

	text := fmt.Sprintf("Подписка %s истекает через %s", subscriptionName, expiresIn)
	params := url.Values{}
	params.Add("text", text)

	err := c.sender.SendTgMessage(params.Encode(), userId)
	if err != nil {
		log.Printf("failed to send tg message: %v", err)
		msg.Nack(false, false)
		return
	}
	msg.Ack(false)
}
