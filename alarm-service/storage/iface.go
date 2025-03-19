package storage

import (
	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
)

// Storage interface for alarms.
// Eventhough the data is stored in memory,
// the design is kept in mind for mograting to a database in future
type AlarmStorage interface {
	SaveAlarm(alarm models.Alarm) error
	GetAlarm(id uuid.UUID) (*models.Alarm, error)
	GetAllAlarms() ([]models.Alarm, error)
	DeleteAlarm(id uuid.UUID) error
	UpdateAlarm(alarm models.Alarm) error
}
