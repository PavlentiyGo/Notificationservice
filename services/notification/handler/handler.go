package notification_handler

import (
	"context"
	"fmt"

	pb "github.com/PavlentiyGo/notification-service/proto"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}
func (h *NotificationHandler) CreateNotification(
	ctx context.Context,
	req *pb.CreateNotificationRequest,
) (*pb.CreateNotificationResponse, error) {

	fmt.Println(req.ChatId, req.Text)

	return &pb.CreateNotificationResponse{Created: false}, nil
}
