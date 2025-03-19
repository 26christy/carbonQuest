package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockACKService is a mock implementation of the ACKService interface
type MockACKService struct {
	mock.Mock
}

func (m *MockACKService) ACKAlarm(alarmID string) error {
	args := m.Called(alarmID)
	return args.Error(0)
}

func (m *MockACKService) GetACKState(alarmID string) (models.ACKState, bool) {
	args := m.Called(alarmID)
	return args.Get(0).(models.ACKState), args.Bool(1)
}

func setupTestRouter() (*gin.Engine, *MockACKService) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockService := new(MockACKService)
	RegisterRoutes(router, mockService)
	return router, mockService
}

func TestACKAlarm_Success(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("ACKAlarm", "123").Return(nil)

	req, _ := http.NewRequest("POST", "/ack/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message":"alarm successfully acknowledged"}`, resp.Body.String())
	mockService.AssertCalled(t, "ACKAlarm", "123")
}

func TestACKAlarm_Failure(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("ACKAlarm", "123").Return(errors.New("failed to ACK alarm"))

	req, _ := http.NewRequest("POST", "/ack/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.JSONEq(t, `{"details":"failed to ACK alarm", "error":"failed to ACK alarm"}`, resp.Body.String())
}

func TestCheckACKState_NotExists(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("GetACKState", "123").Return(models.ACKState{}, false)

	req, _ := http.NewRequest("GET", "/ack/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{
		"alarm_id": "123",
		"acked_at": null,
		"next_notification_at": null,
		"should_notify": true
	}`, resp.Body.String())
}

func TestCheckACKState_Exists_ShouldNotify(t *testing.T) {
	router, mockService := setupTestRouter()

	ackState := models.ACKState{
		AlarmID:            "123",
		ACKedAt:            time.Now().Add(-time.Hour),
		NextNotificationAt: time.Now().Add(-time.Minute), // Time in the past
	}

	mockService.On("GetACKState", "123").Return(ackState, true)

	req, _ := http.NewRequest("GET", "/ack/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"should_notify":true`)
	assert.Contains(t, resp.Body.String(), `"alarm_id":"123"`)
}

func TestCheckACKState_Exists_ShouldNotNotify(t *testing.T) {
	router, mockService := setupTestRouter()

	ackState := models.ACKState{
		AlarmID:            "123",
		ACKedAt:            time.Now().Add(-time.Hour),
		NextNotificationAt: time.Now().Add(time.Minute), // Time in the future
	}

	mockService.On("GetACKState", "123").Return(ackState, true)

	req, _ := http.NewRequest("GET", "/ack/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"should_notify":false`)
	assert.Contains(t, resp.Body.String(), `"alarm_id":"123"`)
}
