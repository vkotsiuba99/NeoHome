package handler

import "github.com/vkotsiuba99/NeoHome/back/internal/domain"

type RespConv struct{}

func NewRespConv() *RespConv {
	return &RespConv{}
}

func (c *RespConv) ToUser(user domain.User) UserResp {
	return UserResp{
		UserID:    user.UserID,
		Email:     user.Email,
		Phone:     user.Phone,
		Login:     user.Login,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (c *RespConv) ToAuth(session domain.Session) AuthResp {
	return AuthResp{
		AccessToken: session.AccessToken,
		ExpiresAt:   session.ExpiresAt,
		User:        c.ToUser(session.User),
	}
}

func (c *RespConv) ToUserWrap(user domain.User) UserWrapResp {
	return UserWrapResp{
		User: c.ToUser(user),
	}
}

func (c *RespConv) ToDevice(device domain.Device) DeviceResp {
	return DeviceResp{
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

func (c *RespConv) ToDeviceWrap(device domain.Device) DeviceWrapResp {
	return DeviceWrapResp{
		Device: c.ToDevice(device),
	}
}

func (c *RespConv) ToDevices(devices []domain.Device) DevicesResp {
	result := make([]DeviceResp, len(devices))
	for i := range devices {
		result[i] = c.ToDevice(devices[i])
	}
	return DevicesResp{Devices: result}
}

func (c *RespConv) ToThreshold(threshold domain.DeviceThreshold) ThresholdResp {
	return ThresholdResp{
		MetricType:  threshold.MetricType,
		MinValue:    threshold.MinValue,
		HasMinValue: threshold.HasMinValue,
		MaxValue:    threshold.MaxValue,
		HasMaxValue: threshold.HasMaxValue,
		Severity:    threshold.Severity,
		UpdatedAt:   threshold.UpdatedAt,
	}
}

func (c *RespConv) ToThresholds(thresholds []domain.DeviceThreshold) ThresholdsResp {
	result := make([]ThresholdResp, len(thresholds))
	for i := range thresholds {
		result[i] = c.ToThreshold(thresholds[i])
	}
	return ThresholdsResp{Thresholds: result}
}

func (c *RespConv) ToTelemetry(telemetry domain.Telemetry) TelemetryResp {
	return TelemetryResp{
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

func (c *RespConv) ToTelemetryList(telemetryRows []domain.Telemetry) TelemetryListResp {
	result := make([]TelemetryResp, len(telemetryRows))
	for i := range telemetryRows {
		result[i] = c.ToTelemetry(telemetryRows[i])
	}
	return TelemetryListResp{Telemetry: result}
}

func (c *RespConv) ToAlert(alert domain.Alert) AlertResp {
	return AlertResp{
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

func (c *RespConv) ToAlerts(alerts []domain.Alert) AlertsResp {
	result := make([]AlertResp, len(alerts))
	for i := range alerts {
		result[i] = c.ToAlert(alerts[i])
	}
	return AlertsResp{Alerts: result}
}

func (c *RespConv) ToAlertWrap(alert domain.Alert) AlertWrapResp {
	return AlertWrapResp{
		Alert: c.ToAlert(alert),
	}
}

func (c *RespConv) ToTelemetryIngest(result domain.IngestResult) TelemetryIngestResp {
	alerts := make([]AlertResp, len(result.Alerts))
	for i := range result.Alerts {
		alerts[i] = c.ToAlert(result.Alerts[i])
	}
	return TelemetryIngestResp{
		Telemetry: c.ToTelemetry(result.Telemetry),
		Alerts:    alerts,
	}
}
