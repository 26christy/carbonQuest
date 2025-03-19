package handlers

import (
	"github.com/26christy/CarbonQuest/notification-service/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, notificationService service.NotificationService) {
	handler := NewNotificationHandler(notificationService)

	notifyGroup := router.Group("/notify")
	{
		notifyGroup.POST("/register-notifier", handler.RegisterNotifier)
		notifyGroup.POST("/", handler.NotificationHandler)
	}
}
