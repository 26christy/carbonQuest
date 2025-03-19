package service

import "github.com/26christy/CarbonQuest/common/models"

type ACKService interface {
	ACKAlarm(alarmID string) error
	GetACKState(alarmID string) (models.ACKState, bool)
}
