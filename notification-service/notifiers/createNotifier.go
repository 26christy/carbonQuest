package notifiers

import "errors"

// CreateNotifier is a factory function that returns a Notifier instance based on type
// Add the notification type as per the requirement
func CreateNotifier(notifierType, param string) (Notifier, error) {
	switch notifierType {
	case "webhook":
		if param == "" {
			return nil, errors.New("webhook URL is missing")
		}
		return NewWebHookNotifier(param), nil
	case "log":
		// LogNotifier does not require a param
		return NewLogNotifier(), nil
	default:
		return nil, errors.New("unsupported notifier type")
	}
}
