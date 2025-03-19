package service

import (
	"github.com/26christy/CarbonQuest/alarm-service/storage"
	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
)

type AlarmService struct {
	store storage.AlarmStorage
}

// Constructor
func NewAlarmService(store storage.AlarmStorage) *AlarmService {
	return &AlarmService{
		store: store,
	}
}

// Create Alarm
func (s *AlarmService) CreateAlarm(alarm models.Alarm) error {
	err := s.store.SaveAlarm(alarm)
	if err != nil {
		return err
	}

	return nil
}

// Get an alarm
func (s *AlarmService) GetAlarm(id uuid.UUID) (*models.Alarm, error) {
	return s.store.GetAlarm(id)
}

// Get all alarms
func (s *AlarmService) GetAllAlarm() ([]models.Alarm, error) {
	return s.store.GetAllAlarms()
}

// Delete an alarm
func (s *AlarmService) DeleteAlarm(id uuid.UUID) error {
	return s.store.DeleteAlarm(id)
}

// Update an alarm
func (s *AlarmService) UpdateAlarm(alarm models.Alarm) error {
	return s.store.UpdateAlarm(alarm)
}
