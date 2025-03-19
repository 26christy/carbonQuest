package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/26christy/CarbonQuest/alarm-service/service"
	"github.com/26christy/CarbonQuest/common/models"
	"github.com/26christy/CarbonQuest/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlarmHandler struct {
	service service.AlarmServiceInterface
}

func NewAlarmHandler(service service.AlarmServiceInterface) *AlarmHandler {
	return &AlarmHandler{
		service: service,
	}
}

func (h *AlarmHandler) createAlarm(c *gin.Context) {
	var req models.Alarm

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request body validation failed",
			"details": err.Error(),
		})
		return
	}

	req.ID = uuid.New()

	// Ensure timestamps are set
	req.CreatedAt = time.Now()
	req.UpdatedAt = req.CreatedAt

	// Set a default status if not provided
	if req.Status == "" {
		req.Status = "triggered" // Default status when an alarm is created
	}

	// Call service to create alarm
	if err := h.service.CreateAlarm(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create alarm",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *AlarmHandler) getAlarm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid alarm id",
			"details": err.Error(),
		})
		return
	}

	alarm, err := h.service.GetAlarm(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "failed to fetch the details for alarm ID: " + id.String(),
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, alarm)
}

func (h *AlarmHandler) getAllAlarm(c *gin.Context) {
	alarms, err := h.service.GetAllAlarm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch the alarms",
			"details": err.Error()})
		return
	}

	if alarms == nil {
		alarms = []models.Alarm{}
	}

	c.JSON(http.StatusOK, gin.H{
		"alarms": alarms,
	})
}

func (h *AlarmHandler) deleteAlarm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid alarm id",
			"details": err.Error(),
		})
		return
	}

	err = h.service.DeleteAlarm(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "failed to delete alarm for alarm ID: " + id.String(),
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *AlarmHandler) updateAlarm(c *gin.Context) {
	id, err := h.parseAlarmID(c)
	if err != nil {
		return
	}

	var req models.UpdateAlarm
	if err := h.bindAndValidateRequest(c, &req); err != nil {
		return
	}

	existingAlarm, err := h.service.GetAlarm(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   fmt.Sprintf("failed to fetch the details for alarm ID: %s", id.String()),
			"details": err.Error(),
		})
		return
	}

	if !h.isValidStateTransition(existingAlarm.Status, req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state transition",
		})
		return
	}

	// Fill missing fields with existing values
	updatedAlarm := h.fillMissingFields(req, *existingAlarm)

	if err := h.service.UpdateAlarm(updatedAlarm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("failed to update the alarm for ID: %s", id.String()),
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedAlarm)
}

// Parses and validates the alarm ID from the request
func (h *AlarmHandler) parseAlarmID(c *gin.Context) (uuid.UUID, error) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid alarm id",
			"details": err.Error(),
		})
	}
	return id, err
}

// Binds and validates the update request body
func (h *AlarmHandler) bindAndValidateRequest(c *gin.Context, req *models.UpdateAlarm) error {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return err
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request body validation failed",
			"details": err.Error(),
		})
		return err
	}

	return nil
}

// Checks if the state transition is valid
func (h *AlarmHandler) isValidStateTransition(currentStatus, newState string) bool {
	validTransitions := map[string][]string{
		"triggered": {"active", "ACK"},
		"active":    {"ACK"},
		"ACK":       {"active"},
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	return exists && contains(allowedTransitions, newState)
}

// Fills missing fields from the existing alarm
func (h *AlarmHandler) fillMissingFields(req models.UpdateAlarm, existing models.Alarm) models.Alarm {
	// If fields are empty, use the existing values
	if req.Name == "" {
		req.Name = existing.Name
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = existing.Timestamp
	}
	if req.Status == "" {
		req.Status = existing.Status
	}

	return models.Alarm{
		ID:        existing.ID,
		Name:      req.Name,
		Timestamp: req.Timestamp,
		Status:    req.Status,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}
}

// Helper function to check if a state transition is allowed
func contains(states []string, target string) bool {
	for _, state := range states {
		if state == target {
			return true
		}
	}
	return false
}
