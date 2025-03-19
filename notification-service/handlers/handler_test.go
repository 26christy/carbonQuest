package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
	"github.com/26christy/CarbonQuest/notification-service/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(alarm models.AlarmEvent) {
	m.Called(alarm)
}

func (m *MockNotificationService) RegisterNotifier(n notifiers.Notifier) {
	m.Called(n)
}

func (m *MockNotificationService) StartNotificationScheduler() {
	m.Called()
}

func (m *MockNotificationService) UpdateACKState(alarmID string, ackTime time.Time) {
	m.Called(alarmID, ackTime)
}

func (m *MockNotificationService) CancelUnACKedNotifications(alarmID string) {
	m.Called(alarmID)
}

func (m *MockNotificationService) ScheduleNextACKedNotification(alarmID string, nextTime time.Time) {
	m.Called(alarmID, nextTime)
}

// Setup helper function for the test router
func setupTestRouter(service service.NotificationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewNotificationHandler(service)
	router := gin.Default()

	router.POST("/notify/register-notifier", handler.RegisterNotifier)
	router.POST("/notify/event", handler.NotificationHandler)

	return router
}

func TestRegisterNotifier_Success(t *testing.T) {
	mockService := new(MockNotificationService)
	mockService.On("RegisterNotifier", mock.Anything).Return(nil)

	router := setupTestRouter(mockService)

	reqBody := `{"type": "log", "param": ""}`
	req, _ := http.NewRequest("POST", "/notify/register-notifier", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertCalled(t, "RegisterNotifier", mock.Anything)
}

func TestRegisterNotifier_MissingFields(t *testing.T) {
	mockService := new(MockNotificationService)
	router := setupTestRouter(mockService)

	// Missing "type" field
	reqBody := `{"param": "http://example.com"}`
	req, _ := http.NewRequest("POST", "/notify/register-notifier", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestRegisterNotifier_InvalidType(t *testing.T) {
	mockService := new(MockNotificationService)
	router := setupTestRouter(mockService)

	reqBody := `{"type": "invalid_type", "param": ""}`
	req, _ := http.NewRequest("POST", "/notify/register-notifier", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestNotificationHandler_Success(t *testing.T) {
	mockService := new(MockNotificationService)
	mockService.On("SendNotification", mock.Anything).Return()

	router := setupTestRouter(mockService)

	reqBody := `{
	"alarm_id": "123",
	"name": "Test Alarm",
	"type": "Test Type",
	"timestamp": "2025-03-19T10:00:00Z"
}`
	req, _ := http.NewRequest("POST", "/notify/event", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertCalled(t, "SendNotification", mock.Anything)
}

func TestNotificationHandler_InvalidRequest(t *testing.T) {
	mockService := new(MockNotificationService)
	router := setupTestRouter(mockService)

	// Invalid request (missing required fields)
	reqBody := `{"invalid_field": "123"}`
	req, _ := http.NewRequest("POST", "/notify/event", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
