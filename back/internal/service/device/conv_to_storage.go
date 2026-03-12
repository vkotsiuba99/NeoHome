package device

import (
	"math"
	"strings"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToStore struct{}

func NewStoreConv() *ConvToStore {
	return &ConvToStore{}
}

func (c *ConvToStore) DomainDeviceToStorage(device domain.Device) storage.Device {
	return storage.Device{
		DeviceID:       device.DeviceID,
		UserID:         device.UserID,
		DeviceName:     device.DeviceName,
		DeviceType:     device.DeviceType,
		RoomName:       device.RoomName,
		LocationID:     device.LocationID,
		LocationName:   device.LocationName,
		Status:         device.Status,
		LastSeenAt:     device.LastSeenAt,
		LastMetricAt:   device.LastMetricAt,
		BatteryLevel:   device.BatteryLevel,
		SignalStrength: device.SignalStrength,
		AddedAt:        device.AddedAt,
		UpdatedAt:      device.UpdatedAt,
	}
}

func (c *ConvToStore) PatchToStorageThreshold(deviceID int64, patch domain.ThresholdPatch, metricType string, severity string, now int64) storage.DeviceThreshold {
	minValue := patch.MinValue
	if !patch.HasMinValue {
		minValue = math.NaN()
	}
	maxValue := patch.MaxValue
	if !patch.HasMaxValue {
		maxValue = math.NaN()
	}

	return storage.DeviceThreshold{
		DeviceID:    deviceID,
		MetricType:  metricType,
		MinValue:    minValue,
		HasMinValue: patch.HasMinValue,
		MaxValue:    maxValue,
		HasMaxValue: patch.HasMaxValue,
		Severity:    severity,
		UpdatedAt:   now,
	}
}

func (c *ConvToStore) PatchesToStorageThresholds(deviceID int64, patches []domain.ThresholdPatch, defaultSeverity string, now int64) ([]storage.DeviceThreshold, bool) {
	rows := make([]storage.DeviceThreshold, 0, len(patches))
	for _, item := range patches {
		metricType := strings.ToLower(strings.TrimSpace(item.MetricType))
		if len(metricType) == 0 || (!item.HasMinValue && !item.HasMaxValue) || (item.HasMinValue && item.HasMaxValue && item.MinValue > item.MaxValue) {
			return nil, false
		}

		severity := strings.ToLower(strings.TrimSpace(item.Severity))
		if len(severity) == 0 {
			severity = defaultSeverity
		}

		rows = append(rows, c.PatchToStorageThreshold(deviceID, item, metricType, severity, now))
	}

	return rows, true
}
