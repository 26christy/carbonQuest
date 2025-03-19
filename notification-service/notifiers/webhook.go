package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/26christy/CarbonQuest/common/models"
)

// WebHookNotifier sends notifications to an external WebHook
type WebHookNotifier struct {
	URL string
}

// NewWebHookNotifier initializes a WebHookNotifier
func NewWebHookNotifier(url string) *WebHookNotifier {
	return &WebHookNotifier{URL: url}
}

// Notify sends an HTTP POST request to the WebHook
func (w *WebHookNotifier) Notify(alarm models.AlarmEvent) error {
	payload, err := json.Marshal(alarm)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WebHookNotifier: Failed to send notification, status code: %d", resp.StatusCode)
	}

	fmt.Printf("WebHook Notification sent to %s\n", w.URL)
	return nil
}
