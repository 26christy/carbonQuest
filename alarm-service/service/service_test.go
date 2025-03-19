package service

import (
	"errors"
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Storage
type MockAlarmStorage struct {
	mock.Mock
}

func (m *MockAlarmStorage) SaveAlarm(alarm models.Alarm) error {
	args := m.Called(alarm)
	return args.Error(0)
}

func (m *MockAlarmStorage) GetAlarm(id uuid.UUID) (*models.Alarm, error) {
	args := m.Called(id)
	if alarm, ok := args.Get(0).(*models.Alarm); ok {
		return alarm, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlarmStorage) GetAllAlarms() ([]models.Alarm, error) {
	args := m.Called()
	if alarms, ok := args.Get(0).([]models.Alarm); ok {
		return alarms, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlarmStorage) DeleteAlarm(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAlarmStorage) UpdateAlarm(alarm models.Alarm) error {
	args := m.Called(alarm)
	return args.Error(0)
}

// Mock Publisher
type MockEventPublisher struct {
	mock.Mock
}

func createTestService() (*AlarmService, *MockAlarmStorage, *MockEventPublisher) {
	mockStorage := new(MockAlarmStorage)
	mockPublisher := new(MockEventPublisher)

	// Ensure all methods of EventPublisher are mocked
	mockPublisher.On("Subscribe", mock.Anything, mock.Anything).Return()

	service := NewAlarmService(mockStorage)
	return service, mockStorage, mockPublisher
}

func TestCreateAlarm(t *testing.T) {
	service, mockStorage, mockPublisher := createTestService()

	alarm := models.Alarm{
		ID:        uuid.New(),
		Name:      "Test Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	mockStorage.On("SaveAlarm", alarm).Return(nil).Once()

	mockPublisher.On("Publish", mock.Anything).Once()

	err := service.CreateAlarm(alarm)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestCreateAlarm_Error(t *testing.T) {
	service, mockStorage, mockPublisher := createTestService()

	alarm := models.Alarm{
		ID:        uuid.New(),
		Name:      "Test Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	mockStorage.On("SaveAlarm", alarm).Return(errors.New("storage error")).Once()

	err := service.CreateAlarm(alarm)
	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())

	// Ensure publisher is not called due to storage failure
	mockPublisher.AssertNotCalled(t, "Publish", mock.Anything)
	mockStorage.AssertExpectations(t)
}

func TestGetAlarm(t *testing.T) {
	service, mockStorage, _ := createTestService()

	alarmID := uuid.New()
	expectedAlarm := &models.Alarm{
		ID:        alarmID,
		Name:      "Sample Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	mockStorage.On("GetAlarm", alarmID).Return(expectedAlarm, nil).Once()

	alarm, err := service.GetAlarm(alarmID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAlarm, alarm)

	mockStorage.AssertExpectations(t)
}

func TestGetAlarm_NotFound(t *testing.T) {
	service, mockStorage, _ := createTestService()

	alarmID := uuid.New()
	mockStorage.On("GetAlarm", alarmID).Return(nil, errors.New("alarm not found")).Once()

	alarm, err := service.GetAlarm(alarmID)
	assert.Error(t, err)
	assert.Nil(t, alarm)
	assert.Equal(t, "alarm not found", err.Error())

	mockStorage.AssertExpectations(t)
}

func TestGetAllAlarm(t *testing.T) {
	service, mockStorage, _ := createTestService()

	expectedAlarms := []models.Alarm{
		{
			ID:        uuid.New(),
			Name:      "Alarm 1",
			Timestamp: time.Now(),
			Status:    "active",
		},
		{
			ID:        uuid.New(),
			Name:      "Alarm 2",
			Timestamp: time.Now(),
			Status:    "inactive",
		},
	}

	mockStorage.On("GetAllAlarms").Return(expectedAlarms, nil).Once()

	alarms, err := service.GetAllAlarm()
	assert.NoError(t, err)
	assert.Equal(t, expectedAlarms, alarms)

	mockStorage.AssertExpectations(t)
}

func TestDeleteAlarm(t *testing.T) {
	service, mockStorage, _ := createTestService()

	alarmID := uuid.New()

	mockStorage.On("DeleteAlarm", alarmID).Return(nil).Once()

	err := service.DeleteAlarm(alarmID)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestUpdateAlarm(t *testing.T) {
	service, mockStorage, _ := createTestService()

	alarm := models.Alarm{
		ID:        uuid.New(),
		Name:      "Updated Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	mockStorage.On("UpdateAlarm", alarm).Return(nil).Once()

	err := service.UpdateAlarm(alarm)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}
