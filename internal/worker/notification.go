package worker

import (
	"context"

	"github.com/truongnqse05461/ewallet/internal/log"
	"github.com/truongnqse05461/ewallet/internal/model"
)

type NotificationWorker struct {
	notificationChannel chan model.Notification
	logger              log.Logger
}

func NewNotificationWorker(logger log.Logger, bufferSize int) *NotificationWorker {
	return &NotificationWorker{
		notificationChannel: make(chan model.Notification, bufferSize),
		logger:              logger,
	}
}

func (nw *NotificationWorker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case notification := <-nw.notificationChannel:
				nw.processNotification(notification)
			case <-ctx.Done():
				nw.logger.Info("notification worker stopped")
				return
			}
		}
	}()
}

func (nw *NotificationWorker) processNotification(notification model.Notification) {
	nw.logger.Infof("[Notification] new notification sent to %s: %s", notification.UserID, notification.Message)
}

func (nw *NotificationWorker) SendNotification(notification model.Notification) {
	nw.notificationChannel <- notification
}
