package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()
	assert.NotNil(t, storage, "NewMemoryStorage should return a non-nil instance")
}

func TestACKAlarm(t *testing.T) {
	storage := NewMemoryStorage()
	alarmID := "test-alarm"

	err := storage.ACKAlarm(alarmID)
	assert.NoError(t, err, "ACKAlarm should not return an error")

	ackState, exists := storage.GetACKState(alarmID)
	assert.True(t, exists, "ACK state should exist after acknowledgment")
	assert.Equal(t, alarmID, ackState.AlarmID, "Alarm ID should match the acknowledged one")

	// Check if ACK time is set correctly
	assert.WithinDuration(t, time.Now(), ackState.ACKedAt, time.Second, "ACKedAt should be close to the current time")

	// Check if NextNotificationAt is 24 hours from ACKedAt
	expectedNextNotification := ackState.ACKedAt.Add(24 * time.Hour)
	assert.WithinDuration(t, expectedNextNotification, ackState.NextNotificationAt, time.Second, "NextNotificationAt should be 24 hours from ACKedAt")
}

func TestGetACKState_NotExists(t *testing.T) {
	storage := NewMemoryStorage()
	alarmID := "non-existent-alarm"

	_, exists := storage.GetACKState(alarmID)
	assert.False(t, exists, "GetACKState should return false for non-existent alarm")
}
