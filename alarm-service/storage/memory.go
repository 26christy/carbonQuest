package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	alarms map[string]models.Alarm
}

// constructor
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		alarms: make(map[string]models.Alarm),
	}
}

// save alarm
func (m *MemoryStorage) SaveAlarm(alarm models.Alarm) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alarms[alarm.ID.String()] = alarm
	return nil
}

// Get an alarm
func (m *MemoryStorage) GetAlarm(id uuid.UUID) (*models.Alarm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	alarm, exists := m.alarms[id.String()]
	if !exists {
		return nil, errors.New("alarm not found")
	}
	return &alarm, nil
}

// Get all alarms
func (m *MemoryStorage) GetAllAlarms() ([]models.Alarm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []models.Alarm
	for _, alarm := range m.alarms {
		result = append(result, alarm)
	}
	return result, nil
}

// Delete an alarm
func (m *MemoryStorage) DeleteAlarm(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.alarms[id.String()]; !exists {
		return errors.New("alarm not found")
	}
	delete(m.alarms, id.String())
	return nil
}

// Update an alarm
func (m *MemoryStorage) UpdateAlarm(alarm models.Alarm) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the alarm exists
	existingAlarm, exists := m.alarms[alarm.ID.String()]
	if !exists {
		return errors.New("alarm not found")
	}

	// Keep CreatedAt the same, update UpdatedAt
	alarm.CreatedAt = existingAlarm.CreatedAt
	alarm.UpdatedAt = time.Now()

	// Save the updated alarm
	m.alarms[alarm.ID.String()] = alarm
	return nil
}
