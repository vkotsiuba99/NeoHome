package telemetry

import (
	"crypto/rand"
	"math"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToDamain struct{}

func NewDomainConv() *ConvToDamain {
	return &ConvToDamain{}
}

func (c *ConvToDamain) StorageDeviceToDomain(dto storage.Device) domain.Device {
	return domain.Device{
		DeviceID:       dto.DeviceID,
		UserID:         dto.UserID,
		DeviceName:     dto.DeviceName,
		DeviceType:     dto.DeviceType,
		RoomName:       dto.RoomName,
		LocationID:     dto.LocationID,
		LocationName:   dto.LocationName,
		Status:         dto.Status,
		LastSeenAt:     dto.LastSeenAt,
		LastMetricAt:   dto.LastMetricAt,
		BatteryLevel:   dto.BatteryLevel,
		SignalStrength: dto.SignalStrength,
		AddedAt:        dto.AddedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

func (c *ConvToDamain) StorageTelemetryToDomain(dto storage.Telemetry) domain.Telemetry {
	return domain.Telemetry{
		TelemetryID:    dto.TelemetryID,
		DeviceID:       dto.DeviceID,
		RecordedAt:     dto.RecordedAt,
		MetricType:     dto.MetricType,
		MetricValue:    dto.MetricValue,
		Unit:           dto.Unit,
		RoomName:       dto.RoomName,
		LocationName:   dto.LocationName,
		BatteryLevel:   dto.BatteryLevel,
		SignalStrength: dto.SignalStrength,
	}
}

func (c *ConvToDamain) IngestToDomainNormalize(cmd domain.TelemetryIngest) (domain.TelemetryIngest, bool) {
	normalized := domain.TelemetryIngest{
		DeviceID:       cmd.DeviceID,
		RecordedAt:     cmd.RecordedAt,
		HasRecordedAt:  cmd.HasRecordedAt,
		MetricType:     strings.ToLower(strings.TrimSpace(cmd.MetricType)),
		MetricValue:    cmd.MetricValue,
		Unit:           strings.TrimSpace(cmd.Unit),
		RoomName:       strings.TrimSpace(cmd.RoomName),
		LocationName:   strings.TrimSpace(cmd.LocationName),
		BatteryLevel:   cmd.BatteryLevel,
		HasBattery:     cmd.HasBattery,
		SignalStrength: cmd.SignalStrength,
		HasSignal:      cmd.HasSignal,
	}

	if normalized.DeviceID <= 0 || len(normalized.MetricType) == 0 {
		return domain.TelemetryIngest{}, false
	}
	if normalized.HasRecordedAt && normalized.RecordedAt <= 0 {
		return domain.TelemetryIngest{}, false
	}

	return normalized, true
}

func (c *ConvToDamain) StorageTelemetryListToDomain(rows []storage.Telemetry) []domain.Telemetry {
	items := make([]domain.Telemetry, 0, len(rows))
	for _, row := range rows {
		items = append(items, c.StorageTelemetryToDomain(row))
	}
	return items
}

func (c *ConvToDamain) IngestToDomainTelemetryAndDevice(cmd domain.TelemetryIngest, device domain.Device, activeStatus string) (domain.Telemetry, domain.Device, error) {
	recordedAt := time.Now().UTC().UnixMilli()
	if cmd.HasRecordedAt {
		recordedAt = cmd.RecordedAt
	}

	roomName := cmd.RoomName
	if len(roomName) == 0 {
		roomName = device.RoomName
	}
	locationName := cmd.LocationName
	if len(locationName) == 0 {
		locationName = device.LocationName
	}

	batteryLevel := device.BatteryLevel
	if cmd.HasBattery {
		batteryLevel = cmd.BatteryLevel
	}
	signalStrength := device.SignalStrength
	if cmd.HasSignal {
		signalStrength = cmd.SignalStrength
	}

	telemetryID := int64(0)
	for telemetryID <= 0 {
		rawID, err := rand.Int(rand.Reader, big.NewInt(1<<53))
		if err != nil {
			return domain.Telemetry{}, domain.Device{}, err
		}
		telemetryID = rawID.Int64()
	}

	telemetry := c.IngestToDomainTelemetry(cmd, device, telemetryID, recordedAt, cmd.MetricType, cmd.Unit, roomName, locationName, batteryLevel, signalStrength)

	updatedDevice := device
	updatedDevice.Status = activeStatus
	updatedDevice.LastSeenAt = recordedAt
	updatedDevice.LastMetricAt = recordedAt
	if cmd.HasBattery {
		updatedDevice.BatteryLevel = cmd.BatteryLevel
	}
	if cmd.HasSignal {
		updatedDevice.SignalStrength = cmd.SignalStrength
	}
	updatedDevice.UpdatedAt = time.Now().UTC().UnixMilli()

	return telemetry, updatedDevice, nil
}

func (c *ConvToDamain) StorageTelemetryLatestByMetric(rows []storage.Telemetry) []storage.Telemetry {
	latestByMetric := make(map[string]storage.Telemetry, len(rows))
	for _, row := range rows {
		metric := strings.ToLower(strings.TrimSpace(row.MetricType))
		current, ok := latestByMetric[metric]
		if !ok || row.RecordedAt > current.RecordedAt || (row.RecordedAt == current.RecordedAt && row.TelemetryID > current.TelemetryID) {
			latestByMetric[metric] = row
		}
	}

	latestRows := make([]storage.Telemetry, 0, len(latestByMetric))
	for _, row := range latestByMetric {
		latestRows = append(latestRows, row)
	}
	sort.Slice(latestRows, func(i, j int) bool {
		return latestRows[i].MetricType < latestRows[j].MetricType
	})

	return latestRows
}

func (c *ConvToDamain) StorageThresholdToDomain(dto storage.DeviceThreshold) domain.DeviceThreshold {
	hasMinValue := dto.HasMinValue || !math.IsNaN(dto.MinValue)
	hasMaxValue := dto.HasMaxValue || !math.IsNaN(dto.MaxValue)

	return domain.DeviceThreshold{
		DeviceID:    dto.DeviceID,
		MetricType:  dto.MetricType,
		MinValue:    dto.MinValue,
		HasMinValue: hasMinValue,
		MaxValue:    dto.MaxValue,
		HasMaxValue: hasMaxValue,
		Severity:    dto.Severity,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func (c *ConvToDamain) IngestToDomainTelemetry(cmd domain.TelemetryIngest, device domain.Device, telemetryID int64, recordedAt int64, metricType string, unit string, roomName string, locationName string, batteryLevel int64, signalStrength int64) domain.Telemetry {
	return domain.Telemetry{
		TelemetryID:    telemetryID,
		DeviceID:       device.DeviceID,
		RecordedAt:     recordedAt,
		MetricType:     metricType,
		MetricValue:    cmd.MetricValue,
		Unit:           unit,
		RoomName:       roomName,
		LocationName:   locationName,
		BatteryLevel:   batteryLevel,
		SignalStrength: signalStrength,
	}
}

func (c *ConvToDamain) TelemetryToDomainAlert(alertID int64, device domain.Device, telemetry domain.Telemetry, severity string, createdAt int64, message string) domain.Alert {
	return domain.Alert{
		AlertID:    alertID,
		LocationID: device.LocationID,
		DeviceID:   telemetry.DeviceID,
		CreatedAt:  createdAt,
		Severity:   severity,
		Message:    message,
	}
}
