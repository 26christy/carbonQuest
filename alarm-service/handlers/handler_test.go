package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the AlarmServiceInterface
type MockAlarmService struct {
	mock.Mock
}

func (m *MockAlarmService) CreateAlarm(alarm models.Alarm) error {
	args := m.Called(alarm)
	return args.Error(0)
}

func (m *MockAlarmService) GetAlarm(id uuid.UUID) (*models.Alarm, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Alarm), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlarmService) GetAllAlarm() ([]models.Alarm, error) {
	args := m.Called()
	return args.Get(0).([]models.Alarm), args.Error(1)
}

func (m *MockAlarmService) UpdateAlarm(alarm models.Alarm) error {
	args := m.Called(alarm)
	return args.Error(0)
}

func (m *MockAlarmService) DeleteAlarm(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter(handler *AlarmHandler) *gin.Engine {
	router := gin.Default()
	router.POST("/alarms", handler.createAlarm)
	router.GET("/alarms/:id", handler.getAlarm)
	router.GET("/alarms", handler.getAllAlarm)
	router.PUT("/alarms/:id", handler.updateAlarm)
	router.DELETE("/alarms/:id", handler.deleteAlarm)
	return router
}

func TestCreateAlarm(t *testing.T) {
	mockService := new(MockAlarmService)
	handler := NewAlarmHandler(mockService)
	router := setupRouter(handler)

	reqBody := `{"name":"Test Alarm","timestamp":"2023-03-17T15:04:05Z"}`
	req, _ := http.NewRequest(http.MethodPost, "/alarms", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("CreateAlarm", mock.Anything).Return(nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertCalled(t, "CreateAlarm", mock.Anything)
}

func TestGetAlarm(t *testing.T) {
	mockService := new(MockAlarmService)
	handler := NewAlarmHandler(mockService)
	router := setupRouter(handler)

	alarmID := uuid.New()
	alarm := &models.Alarm{
		ID:        alarmID,
		Name:      "Test Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	mockService.On("GetAlarm", alarmID).Return(alarm, nil)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/alarms/%s", alarmID.String()), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var fetchedAlarm models.Alarm
	err := json.Unmarshal(resp.Body.Bytes(), &fetchedAlarm)
	assert.NoError(t, err)
	assert.Equal(t, alarm.ID, fetchedAlarm.ID)
	mockService.AssertCalled(t, "GetAlarm", alarmID)
}

func TestGetAllAlarm(t *testing.T) {
	mockService := new(MockAlarmService)
	handler := NewAlarmHandler(mockService)

	// Set up the Gin router with the handler
	router := gin.Default()
	router.GET("/alarms", handler.getAllAlarm)

	// Test case 1: Successfully fetch all alarms
	alarms := []models.Alarm{
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

	mockService.On("GetAllAlarm").Return(alarms, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/alarms", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var successResponse map[string][]models.Alarm
	err := json.Unmarshal(resp.Body.Bytes(), &successResponse)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(successResponse["alarms"]))
	assert.Equal(t, "Alarm 1", successResponse["alarms"][0].Name)
	assert.Equal(t, "Alarm 2", successResponse["alarms"][1].Name)

	// Test case 2: No alarms found (should return an empty list)
	mockService.On("GetAllAlarm").Return([]models.Alarm{}, nil).Once()

	req, _ = http.NewRequest(http.MethodGet, "/alarms", nil)
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &successResponse)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(successResponse["alarms"]))

	// Ensure the mock expectations are met
	mockService.AssertExpectations(t)
}

func TestUpdateAlarm(t *testing.T) {
	mockService := new(MockAlarmService)
	handler := NewAlarmHandler(mockService)
	router := setupRouter(handler)

	alarmID := uuid.New()
	existingAlarm := &models.Alarm{
		ID:        alarmID,
		Name:      "Old Alarm",
		Timestamp: time.Now(),
		Status:    "active",
	}

	// Mocking the GetAlarm method to return the existing alarm
	mockService.On("GetAlarm", alarmID).Return(existingAlarm, nil)

	// Mocking the UpdateAlarm method to simulate a successful update
	mockService.On("UpdateAlarm", mock.MatchedBy(func(alarm models.Alarm) bool {
		return alarm.ID == alarmID && alarm.Name == "Updated Alarm" && alarm.Status == "ACK"
	})).Return(nil)

	// Preparing a valid JSON request body
	reqBody := `{"name":"Updated Alarm","status":"ACK"}`
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/alarms/%s", alarmID.String()), bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	// Performing the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Asserting the response status
	assert.Equal(t, http.StatusOK, resp.Code, "Expected status 200 OK")

	// Verifying that the expected methods were called
	mockService.AssertCalled(t, "GetAlarm", alarmID)
	mockService.AssertCalled(t, "UpdateAlarm", mock.MatchedBy(func(alarm models.Alarm) bool {
		return alarm.ID == alarmID && alarm.Name == "Updated Alarm" && alarm.Status == "ACK"
	}))
}

func TestDeleteAlarm(t *testing.T) {
	mockService := new(MockAlarmService)
	handler := NewAlarmHandler(mockService)
	router := setupRouter(handler)

	alarmID := uuid.New()
	mockService.On("DeleteAlarm", alarmID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/alarms/%s", alarmID.String()), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockService.AssertCalled(t, "DeleteAlarm", alarmID)
}
