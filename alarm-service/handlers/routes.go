package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, alarmService *AlarmHandler) {
	api := router.Group("/alarms")
	{
		api.POST("/", alarmService.createAlarm)
		api.GET("/:id", alarmService.getAlarm)
		api.GET("/", alarmService.getAllAlarm)
		api.DELETE("/:id", alarmService.deleteAlarm)
		api.PUT("/:id", alarmService.updateAlarm)
	}
}
