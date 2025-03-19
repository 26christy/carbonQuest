package storage

import "github.com/26christy/CarbonQuest/common/models"

type ACKStorage interface {
	ACKAlarm(alarmID string) error
	GetACKState(alarmID string) (models.ACKState, bool)
}
