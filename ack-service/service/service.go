package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/26christy/CarbonQuest/ack-service/storage"
	"github.com/26christy/CarbonQuest/common/models"
)

type ACKServiceImpl struct {
	storage    storage.ACKStorage
	httpClient *http.Client
}

func NewACKService(storage storage.ACKStorage, client *http.Client) ACKService {
	return &ACKServiceImpl{
		storage:    storage,
		httpClient: client,
	}
}

// ACKAlarm marks an alarm as acknowledged and sets next notification time
func (s *ACKServiceImpl) ACKAlarm(alarmID string) error {
	_, success := s.callUpdateAlarm(alarmID)
	if !success {
		return errors.New("failed to update the alarm status")
	}
	return s.storage.ACKAlarm(alarmID)
}

// ShouldNotify checks if the alarm should be notified based on ACK state
func (s *ACKServiceImpl) GetACKState(alarmID string) (models.ACKState, bool) {
	return s.storage.GetACKState(alarmID)
}

func (s *ACKServiceImpl) callUpdateAlarm(alarmID string) (models.UpdateAlarm, bool) {
	host := os.Getenv("HOST")
	port := os.Getenv("ALARM_SERVICE_PORT")
	url := fmt.Sprintf("http://%s:%s/alarms/%s", host, port, alarmID)

	// Create the request payload
	payload := map[string]string{
		"status": "ACK",
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
