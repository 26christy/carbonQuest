package notifiers

import "github.com/26christy/CarbonQuest/common/models"

// Generic interface for integrating any new notification system.
// All notifiers should implement a common interface,
// making the service agnostic to the type of notification.
type Notifier interface {
	Notify(alarm models.AlarmEvent) error
}
