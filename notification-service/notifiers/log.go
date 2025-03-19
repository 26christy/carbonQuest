package notifiers

import (
	"fmt"

	"github.com/26christy/CarbonQuest/common/models"
)

// LogNotifier logs the alarm event to the console
type LogNotifier struct{}

// NewLogNotifier initializes a LogNotifier
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// Notify logs the alarm event
func (l *LogNotifier) Notify(alarm models.AlarmEvent) error {
	fmt.Printf("Log Notification: Alarm [%s] - Status: %s\n", alarm.Name, alarm.Type)
	return nil
}
