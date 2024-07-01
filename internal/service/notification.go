package service

import (
	"context"

	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/worker"
)

type NotificationService interface {
	SendNotification(ctx context.Context, notification model.Notification) error
}

type notificationSvc struct {
	notificationWorker *worker.NotificationWorker
}

func (n *notificationSvc) SendNotification(ctx context.Context, notification model.Notification) error {
	n.notificationWorker.SendNotification(notification)
	return nil
}

func NewNotificationService(worker *worker.NotificationWorker) NotificationService {
	return &notificationSvc{
		notificationWorker: worker,
	}
}
