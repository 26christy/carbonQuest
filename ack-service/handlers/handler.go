package handlers

import (
	"net/http"
	"time"

	"github.com/26christy/CarbonQuest/ack-service/service"
	"github.com/gin-gonic/gin"
)

type ACKHandler struct {
	service service.ACKService
}

func NewACKHandler(service service.ACKService) *ACKHandler {
	return &ACKHandler{
		service: service,
	}
}

// ACKAlarm handles the ACK request
func (h *ACKHandler) ACKAlarm(c *gin.Context) {
	alarmID := c.Param("id")
	if alarmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alarm id missing"})
		return
	}

	err := h.service.ACKAlarm(alarmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to ACK alarm",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "alarm successfully acknowledged"})
}

func (h *ACKHandler) CheckACKState(c *gin.Context) {
	alarmID := c.Param("id")
	if alarmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alarm id missing"})
		return
	}

	ackState, exists := h.service.GetACKState(alarmID)
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"alarm_id":             alarmID,
			"acked_at":             nil,
			"next_notification_at": nil,
			"should_notify":        true,
		})
		return
	}

	// Otherwise, return actual ACK state
	shouldNotify := time.Now().After(ackState.NextNotificationAt)

	c.JSON(http.StatusOK, gin.H{
		"alarm_id":             ackState.AlarmID,
		"acked_at":             ackState.ACKedAt,
		"next_notification_at": ackState.NextNotificationAt,
		"should_notify":        shouldNotify,
	})
}
