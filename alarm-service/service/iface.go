package service

import (
	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
)

type AlarmServiceInterface interface {
	CreateAlarm(alarm models.Alarm) error
	GetAlarm(id uuid.UUID) (*models.Alarm, error)
	GetAllAlarm() ([]models.Alarm, error)
	DeleteAlarm(id uuid.UUID) error
	UpdateAlarm(alarm models.Alarm) error
}
