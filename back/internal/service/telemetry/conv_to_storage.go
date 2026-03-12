package telemetry

import (
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

func (c *ConvToStore) DomainTelemetryToStorage(telemetry domain.Telemetry) storage.Telemetry {
	return storage.Telemetry{
		TelemetryID:    telemetry.TelemetryID,
		DeviceID:       telemetry.DeviceID,
		RecordedAt:     telemetry.RecordedAt,
		MetricType:     telemetry.MetricType,
		MetricValue:    telemetry.MetricValue,
		Unit:           telemetry.Unit,
		RoomName:       telemetry.RoomName,
		LocationName:   telemetry.LocationName,
		BatteryLevel:   telemetry.BatteryLevel,
		SignalStrength: telemetry.SignalStrength,
	}
}

func (c *ConvToStore) DomainAlertToStorage(alert domain.Alert) storage.Alert {
	return storage.Alert{
		AlertID:    alert.AlertID,
		LocationID: alert.LocationID,
		DeviceID:   alert.DeviceID,
		CreatedAt:  alert.CreatedAt,
		Severity:   alert.Severity,
		Message:    alert.Message,
		IsResolved: alert.IsResolved,
		ResolvedAt: alert.ResolvedAt,
	}
}
