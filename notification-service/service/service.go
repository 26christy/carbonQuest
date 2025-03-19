package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/notification-service/notifiers"
)

// Concrete implementation of NotificationService
type notificationServiceImpl struct {
	notifiers         []notifiers.Notifier
	alarms            map[string]models.AlarmEvent
	notificationState map[string]NotificationState
	httpClient        *http.Client
	mu                sync.Mutex
}

// NewNotificationService initializes the notification service
func NewNotificationService(client *http.Client) NotificationService {
	return &notificationServiceImpl{
		notifiers:         []notifiers.Notifier{},
		alarms:            make(map[string]models.AlarmEvent),
		notificationState: make(map[string]NotificationState),
		httpClient:        client,
	}
}

// SendNotification sends notifications to all registered notifiers
func (s *notificationServiceImpl) SendNotification(alarm models.AlarmEvent) {
	fmt.Printf("[NotificationService] Sending alarm: %+v\n", alarm)

	for _, notifier := range s.notifiers {
		err := notifier.Notify(alarm)
		if err != nil {
			fmt.Printf("[Error] Failed to send notification: %v\n", err)
		} else {
			fmt.Printf("[Success] Notification sent via %T\n", notifier)
		}
	}
}

// RegisterNotifier dynamically adds a new notifier
func (s *notificationServiceImpl) RegisterNotifier(notifier notifiers.Notifier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifiers = append(s.notifiers, notifier)
	fmt.Printf("[NotificationService] Registered a new notifier: %v\n", notifier)
}

func (s *notificationServiceImpl) StartNotificationScheduler() {
	fmt.Println("[DEBUG] Notification Scheduler Started")
	ticker := time.NewTicker(1 * time.Minute) // Check every minute

	go func() {
		for t := range ticker.C {
			fmt.Println("[DEBUG] Running scheduler at", t)
			s.checkAndSendNotifications()
			fmt.Println("[DEBUG] Scheduler completed iteration at", time.Now())
		}
	}()
}

func (s *notificationServiceImpl) checkAndSendNotifications() {
	now := time.Now()
	fmt.Println("[DEBUG] Checking notifications at", now)

	// Fetch latest alarms before checking
	alarms, err := s.fetchAlarmsFromAlarmService()
	if err != nil {
		fmt.Println("[ERROR] Failed to fetch alarms, skipping notification check")
		return
	}

	if len(alarms) == 0 {
		fmt.Println("[DEBUG] No alarms found")
		return
	}

	for _, alarm := range alarms {
		s.processAlarm(alarm, now)
	}
}

// Process each alarm based on its state (ACKed / unACKed)
func (s *notificationServiceImpl) processAlarm(alarm models.AlarmEvent, now time.Time) {
	switch alarm.Type {
	case "triggered":
		if time.Now().After(alarm.Timestamp) {
			fmt.Printf("[INFO] Sending first notification for alarm %s (Triggered)\n", alarm.AlarmID)
			s.SendNotification(alarm)
			s.callUpdateAlarm(alarm.AlarmID, "active")
			s.mu.Lock()
			defer s.mu.Unlock()

			s.notificationState[alarm.AlarmID] = NotificationState{
				FirstNotificationSent: true,
				LastNotificationAt:    now,
			}
		}
	case "active":
		s.handleUnACKedAlarm(s.alarms[alarm.AlarmID], s.notificationState[alarm.AlarmID], now)
	case "ACK":
		s.handleACKedAlarm(s.alarms[alarm.AlarmID], s.notificationState[alarm.AlarmID], now)
	}
}

type NotificationState struct {
	FirstNotificationSent bool
	LastNotificationAt    time.Time
}

// Handle ACKed alarms notifications
func (s *notificationServiceImpl) handleACKedAlarm(alarm models.AlarmEvent, n NotificationState, now time.Time) {
	ackDurationStr := os.Getenv("ACK_DURATION")
	ackDurationInt, err := strconv.Atoi(ackDurationStr)
	if err != nil {
		fmt.Printf("Error parsing ACK_DURATION: %v", err)
		ackDurationInt = 1440 // default is 24 hours
	}
	ackDuration := time.Duration(ackDurationInt) * time.Minute

	if now.After(n.LastNotificationAt.Add(ackDuration)) {
		fmt.Println("[DEBUG] Sending reminder for ACKed alarm:", alarm.AlarmID)
		s.SendNotification(alarm)

		s.mu.Lock()
		defer s.mu.Unlock()

		// Update the notification state
		s.notificationState[alarm.AlarmID] = NotificationState{
			FirstNotificationSent: true,
			LastNotificationAt:    now, // Current time as the last notification sent
		}
	} else {
		fmt.Println("[DEBUG] Skipping ACKed alarm, NextNotificationAt not reached")
	}
}

// Send a reminder to unacked alarm if time duration has met
func (s *notificationServiceImpl) handleUnACKedAlarm(alarm models.AlarmEvent, n NotificationState, now time.Time) {
	ackDurationStr := os.Getenv("UNACK_DURATION")
	ackDurationInt, err := strconv.Atoi(ackDurationStr)
	if err != nil {
		fmt.Printf("Error parsing ACK_DURATION: %v", err)
		ackDurationInt = 120 // default is 2 hours.
	}

	ackDuration := time.Duration(ackDurationInt) * time.Minute
	if now.After(n.LastNotificationAt.Add(ackDuration)) {
		fmt.Println("[DEBUG] Sending reminder for unACKed alarm:", alarm.AlarmID)
		s.SendNotification(alarm)
		s.mu.Lock()
		defer s.mu.Unlock()

		// Update the notification state
		s.notificationState[alarm.AlarmID] = NotificationState{
			FirstNotificationSent: true,
			LastNotificationAt:    now, // Current time as the last notification sent
		}
	} else {
		fmt.Println("[DEBUG] Skipping unACKed alarm, NextNotificationAt not reached")
	}
}

type alarmResponse struct {
	Alarms []models.Alarm `json:"alarms"`
}

// fetchAlarmsFromAlarmService retrieves alarms from alarm-service.
func (s *notificationServiceImpl) fetchAlarmsFromAlarmService() ([]models.AlarmEvent, error) {
	fmt.Println("[DEBUG] Fetching alarms from alarm-service...")

	host := os.Getenv("HOST")
	port := os.Getenv("ALARM_SERVICE_PORT")
	url := fmt.Sprintf("http://%s:%s/alarms", host, port)

	// call the alarm-service api
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch alarms: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ERROR] Unexpected status code from alarm-service: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code from alarm-service: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read response body: %v\n", err)
		return nil, err
	}

	// Decode response into `alarmResponse` struct
	var alarmResp alarmResponse
	if err := json.Unmarshal(body, &alarmResp); err != nil {
		fmt.Printf("[ERROR] Failed to decode alarms: %v\n", err)
		return nil, err
	}

	var alarmEvents []models.AlarmEvent
	for _, alarm := range alarmResp.Alarms {
		alarmEvents = append(alarmEvents, models.AlarmEvent{
			AlarmID:   alarm.ID.String(),
			Name:      alarm.Name,
			Type:      alarm.Status,
			Timestamp: alarm.Timestamp,
		})
	}

	// Store fetched alarms in memory
	s.mu.Lock()
	defer s.mu.Unlock()

	// clear old alarms
	s.alarms = make(map[string]models.AlarmEvent)
	for _, event := range alarmEvents {
		s.alarms[event.AlarmID] = event
	}

	fmt.Printf("[DEBUG] Fetched %d alarms from alarm-service\n", len(alarmResp.Alarms))
	return alarmEvents, nil
}

// updateAlarmStatus updates the status of an alarm in the alarm-service.
func (s *notificationServiceImpl) callUpdateAlarm(alarmID string, status string) (models.UpdateAlarm, bool) {
	host := os.Getenv("HOST")
	port := os.Getenv("ALARM_SERVICE_PORT")
	url := fmt.Sprintf("http://%s:%s/alarms/%s", host, port, alarmID)

	// Create the request payload
	payload := map[string]string{
		"status": status,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.UpdateAlarm{}, false
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return models.UpdateAlarm{}, false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return models.UpdateAlarm{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.UpdateAlarm{}, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.UpdateAlarm{}, false
	}

	var updatedAlarm models.UpdateAlarm
	if err := json.Unmarshal(body, &updatedAlarm); err != nil {
		return models.UpdateAlarm{}, false
	}

	return updatedAlarm, true
}
