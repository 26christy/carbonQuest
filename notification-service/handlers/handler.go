package handlers

import (
	"fmt"
	"net/http"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/common/utils"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
	"github.com/26christy/CarbonQuest/notification-service/service"
	"github.com/gin-gonic/gin"
)

// NotificationHandler handles API requests
type NotificationHandler struct {
	service service.NotificationService
}

// NewNotificationHandler initializes the handler
func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// RegisterNotifier API to register WebHooks
func (h *NotificationHandler) RegisterNotifier(c *gin.Context) {
	var request struct {
		Type  string `json:"type" binding:"required"`
		Param string `json:"param"` // URL for webhook, email for email notifier, etc.
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request body validation failed",
			"details": err.Error(),
		})
		return
	}

	notifier, err := notifiers.CreateNotifier(request.Type, request.Param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to create notifier instance",
			"details": err.Error(),
		})
		return
	}

	h.service.RegisterNotifier(notifier)
	c.JSON(http.StatusOK, gin.H{"message": "Notifier registered successfully"})
}

func (h *NotificationHandler) NotificationHandler(c *gin.Context) {
	var alarmEvent models.AlarmEvent

	if err := c.ShouldBindJSON(&alarmEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(alarmEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request body validation failed",
			"details": err.Error(),
		})
		return
	}

	fmt.Println("[Notification Service] Received event:", alarmEvent)
	h.service.SendNotification(alarmEvent)

	c.JSON(http.StatusOK, gin.H{"message": "Notification received successfully"})
}
