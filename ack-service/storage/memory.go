package storage

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/26christy/CarbonQuest/common/models"
)

type memoryStorage struct {
	mu      sync.RWMutex
	ackData map[string]models.ACKState
}

func NewMemoryStorage() ACKStorage {
	return &memoryStorage{
		ackData: make(map[string]models.ACKState),
	}
}

func (s *memoryStorage) ACKAlarm(alarmID string) error {
	ackDurationStr := os.Getenv("ACK_DURATION")
	ackDurationInt, err := strconv.Atoi(ackDurationStr)
	if err != nil {
		fmt.Printf("Error parsing ACK_DURATION: %v", err)
		ackDurationInt = 1440 // default is 24 hours
	}
	ackDuration := time.Duration(ackDurationInt) * time.Minute
	s.mu.Lock()
	defer s.mu.Unlock()

	ackState := models.ACKState{
		AlarmID:            alarmID,
		ACKedAt:            time.Now(),
		NextNotificationAt: time.Now().Add(ackDuration),
	}

	s.ackData[alarmID] = ackState
	return nil
}

func (s *memoryStorage) GetACKState(alarmID string) (models.ACKState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ackState, exists := s.ackData[alarmID]
	return ackState, exists
}
