package notification_handler

import (
	"context"
	"fmt"

	"github.com/PavlentiyGo/notification-service/proto/notification"
)

type NotificationHandler struct {
	notification.UnimplementedNotificationServiceServer
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}
func (h *NotificationHandler) CreateNotification(
	ctx context.Context,
	req *notification.CreateNotificationRequest,
) (*notification.CreateNotificationResponse, error) {

	fmt.Println(req.ChatId, req.Text)

	return &notification.CreateNotificationResponse{Created: false}, nil
}
