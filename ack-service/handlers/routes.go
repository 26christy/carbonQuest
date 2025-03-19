package handlers

import (
	"github.com/26christy/CarbonQuest/ack-service/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, ackService service.ACKService) {
	handler := NewACKHandler(ackService)

	ackGroup := router.Group("/ack")
	{
		ackGroup.POST("/:id", handler.ACKAlarm)
		ackGroup.GET("/:id", handler.CheckACKState)
	}
}
