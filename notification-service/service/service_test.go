package service

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Notify(alarm models.AlarmEvent) error {
	args := m.Called(alarm)
	return args.Error(0)
}

// MockHTTPClient to simulate HTTP client behavior
type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)), // Mock response body
	}, nil
}

// MockRoundTripper simulates the RoundTripper interface of http.Client
type MockRoundTripper struct{}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)), // Mock response body
	}, nil
}

// setupNotificationService initializes the service with a mock HTTP client
func setupNotificationService() *notificationServiceImpl {
	mockClient := &http.Client{
		Transport: &MockRoundTripper{}, // Use the mock RoundTripper here
	}

	return &notificationServiceImpl{
		notifiers:         []notifiers.Notifier{},
		alarms:            make(map[string]models.AlarmEvent),
		notificationState: make(map[string]NotificationState),
		httpClient:        mockClient,
		mu:                sync.Mutex{},
	}
}

func TestSendNotification(t *testing.T) {
	service := setupNotificationService()
	mockNotifier := new(MockNotifier)
	service.RegisterNotifier(mockNotifier)

	alarm := models.AlarmEvent{
		AlarmID:   "123",
		Name:      "Test Alarm",
		Type:      "triggered",
		Timestamp: time.Now(),
	}

	mockNotifier.On("Notify", alarm).Return(nil)

	service.SendNotification(alarm)
	mockNotifier.AssertCalled(t, "Notify", alarm)
}

func TestRegisterNotifier(t *testing.T) {
	service := setupNotificationService()
	mockNotifier := new(MockNotifier)

	assert.Equal(t, 0, len(service.notifiers))
	service.RegisterNotifier(mockNotifier)
	assert.Equal(t, 1, len(service.notifiers))
}

func TestStartNotificationScheduler(t *testing.T) {
	service := setupNotificationService()

	service.StartNotificationScheduler()
	assert.NotNil(t, service)
}

func TestHandleACKedAlarm(t *testing.T) {
	service := setupNotificationService()

	os.Setenv("ACK_DURATION", "1") // 1 minute
	alarm := models.AlarmEvent{
		AlarmID:   "123",
		Name:      "ACKed Alarm",
		Type:      "ACK",
		Timestamp: time.Now(),
	}
	state := NotificationState{
		LastNotificationAt: time.Now().Add(-2 * time.Minute),
	}

	service.handleACKedAlarm(alarm, state, time.Now())
	assert.True(t, service.notificationState["123"].FirstNotificationSent)
}

func TestHandleUnACKedAlarm(t *testing.T) {
	service := setupNotificationService()

	os.Setenv("UNACK_DURATION", "1") // 1 minute
	alarm := models.AlarmEvent{
		AlarmID:   "123",
		Name:      "UnACKed Alarm",
		Type:      "active",
		Timestamp: time.Now(),
	}
	state := NotificationState{
		LastNotificationAt: time.Now().Add(-2 * time.Minute),
	}

	service.handleUnACKedAlarm(alarm, state, time.Now())
	assert.True(t, service.notificationState["123"].FirstNotificationSent)
}
