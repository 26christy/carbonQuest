package service

import (
	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
)

type NotificationService interface {
	SendNotification(alarm models.AlarmEvent)
	RegisterNotifier(n notifiers.Notifier)
	StartNotificationScheduler()
}
