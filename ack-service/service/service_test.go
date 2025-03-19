package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the ACKStorage interface
type MockACKStorage struct {
	mock.Mock
}

func (m *MockACKStorage) ACKAlarm(alarmID string) error {
	args := m.Called(alarmID)
	return args.Error(0)
}

func (m *MockACKStorage) GetACKState(alarmID string) (models.ACKState, bool) {
	args := m.Called(alarmID)
	return args.Get(0).(models.ACKState), args.Bool(1)
}

// Mock for the HTTP Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func setupTestService() (*ACKServiceImpl, *MockACKStorage, *MockHTTPClient) {
	mockStorage := new(MockACKStorage)
	mockHTTPClient := new(MockHTTPClient)

	// Create an http.Client with the mock HTTP client
	httpClient := &http.Client{
		Transport: mockHTTPClient, // Using the mock as a RoundTripper
	}

	service := &ACKServiceImpl{
		storage:    mockStorage,
		httpClient: httpClient,
	}

	return service, mockStorage, mockHTTPClient
}

func TestACKAlarm_Success(t *testing.T) {
	service, mockStorage, mockHTTPClient := setupTestService()

	alarmID := "test-alarm-id"
	mockStorage.On("ACKAlarm", alarmID).Return(nil)

	// Create a mock response for the HTTP request
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"status":"acknowledged"}`)),
	}

	mockHTTPClient.On("RoundTrip", mock.Anything).Return(mockResponse, nil)

	// Call the method under test
	err := service.ACKAlarm(alarmID)

	// Assertions
	assert.NoError(t, err)
	mockStorage.AssertCalled(t, "ACKAlarm", alarmID)
	mockHTTPClient.AssertCalled(t, "RoundTrip", mock.Anything)
}

func TestGetACKState(t *testing.T) {
	service, mockStorage, _ := setupTestService()

	alarmID := "test-alarm-id"
	expectedACKState := models.ACKState{
		AlarmID:      alarmID,
		ShouldNotify: true,
		// ACKedAt: "2025-03-18T15:24:00+05:30",
	}

	// Test Case 1: ACK state exists
	mockStorage.On("GetACKState", alarmID).Return(expectedACKState, true)

	state, exists := service.GetACKState(alarmID)

	assert.True(t, exists, "ACK state should exist")
	assert.Equal(t, expectedACKState, state, "Returned ACK state should match expected")

	// Test Case 2: ACK state does not exist
	mockStorage.On("GetACKState", "invalid-alarm-id").Return(models.ACKState{}, false)

	state, exists = service.GetACKState("invalid-alarm-id")

	assert.False(t, exists, "ACK state should not exist for an invalid alarm ID")
	assert.Equal(t, models.ACKState{}, state, "Returned ACK state should be empty for an invalid alarm ID")

	mockStorage.AssertExpectations(t)
}
