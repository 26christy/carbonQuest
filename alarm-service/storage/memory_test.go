package storage

import (
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	// Sample alarm for testing
	alarm := models.Alarm{
		ID:        uuid.New(),
		Name:      "Test Alarm",
		Timestamp: time.Now(),
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test SaveAlarm
	err := storage.SaveAlarm(alarm)
	assert.NoError(t, err, "SaveAlarm should not return an error")

	// Test GetAlarm - success
	retrievedAlarm, err := storage.GetAlarm(alarm.ID)
	assert.NoError(t, err, "GetAlarm should not return an error for existing alarm")
	assert.Equal(t, alarm.ID, retrievedAlarm.ID, "Retrieved alarm ID should match the saved one")
	assert.Equal(t, alarm.Name, retrievedAlarm.Name, "Retrieved alarm Name should match the saved one")

	// Test GetAlarm - failure (non-existent alarm)
	_, err = storage.GetAlarm(uuid.New())
	assert.Error(t, err, "GetAlarm should return an error for non-existent alarm")
	assert.Equal(t, "alarm not found", err.Error())

	// Test GetAllAlarms
	alarms, err := storage.GetAllAlarms()
	assert.NoError(t, err, "GetAllAlarms should not return an error")
	assert.Len(t, alarms, 1, "There should be exactly one alarm in the storage")

	// Test UpdateAlarm - success
	updatedAlarm := alarm
	updatedAlarm.Name = "Updated Alarm"
	err = storage.UpdateAlarm(updatedAlarm)
	assert.NoError(t, err, "UpdateAlarm should not return an error")

	retrievedUpdatedAlarm, _ := storage.GetAlarm(alarm.ID)
	assert.Equal(t, "Updated Alarm", retrievedUpdatedAlarm.Name, "Alarm name should be updated")
	assert.NotEqual(t, alarm.UpdatedAt, retrievedUpdatedAlarm.UpdatedAt, "UpdatedAt should be modified")

	// Test UpdateAlarm - failure (non-existent alarm)
	nonExistentAlarm := models.Alarm{
		ID:        uuid.New(),
		Name:      "Non-Existent Alarm",
		Timestamp: time.Now(),
	}
	err = storage.UpdateAlarm(nonExistentAlarm)
	assert.Error(t, err, "UpdateAlarm should return an error for non-existent alarm")
	assert.Equal(t, "alarm not found", err.Error())

	// Test DeleteAlarm - success
	err = storage.DeleteAlarm(alarm.ID)
	assert.NoError(t, err, "DeleteAlarm should not return an error")

	_, err = storage.GetAlarm(alarm.ID)
	assert.Error(t, err, "GetAlarm should return an error after deletion")
	assert.Equal(t, "alarm not found", err.Error())

	// Test DeleteAlarm - failure (non-existent alarm)
	err = storage.DeleteAlarm(alarm.ID)
	assert.Error(t, err, "DeleteAlarm should return an error for non-existent alarm")
	assert.Equal(t, "alarm not found", err.Error())
}
