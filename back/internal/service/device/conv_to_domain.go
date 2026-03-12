package device

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

func (c *ConvToDamain) CreateToDomainNormalize(cmd domain.DeviceCreate) domain.DeviceCreate {
	return domain.DeviceCreate{
		UserID:       cmd.UserID,
		DeviceName:   strings.TrimSpace(cmd.DeviceName),
		DeviceType:   strings.TrimSpace(cmd.DeviceType),
		RoomName:     strings.TrimSpace(cmd.RoomName),
		LocationID:   cmd.LocationID,
		LocationName: strings.TrimSpace(cmd.LocationName),
		Status:       strings.TrimSpace(cmd.Status),
	}
}

func (c *ConvToDamain) CreateToDomainDevice(cmd domain.DeviceCreate, deviceID int64, now int64) domain.Device {
	return domain.Device{
		DeviceID:     deviceID,
		UserID:       cmd.UserID,
		DeviceName:   cmd.DeviceName,
		DeviceType:   cmd.DeviceType,
		RoomName:     cmd.RoomName,
		LocationID:   cmd.LocationID,
		LocationName: cmd.LocationName,
		Status:       cmd.Status,
		AddedAt:      now,
		UpdatedAt:    now,
	}
}

func (c *ConvToDamain) CreateToDomainDeviceAuto(cmd domain.DeviceCreate) (domain.Device, error) {
	nextID := func() (int64, error) {
		for {
			rawID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				return 0, err
			}
			if id := rawID.Int64(); id > 0 {
				return id, nil
			}
		}
	}

	deviceID, err := nextID()
	if err != nil {
		return domain.Device{}, err
	}
	if cmd.LocationID <= 0 {
		locationID, idErr := nextID()
		if idErr != nil {
			return domain.Device{}, idErr
		}
		cmd.LocationID = locationID
	}

	now := time.Now().UTC().UnixMilli()
	return c.CreateToDomainDevice(cmd, deviceID, now), nil
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

func (c *ConvToDamain) StorageDevicesToDomain(rows []storage.Device) []domain.Device {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].AddedAt == rows[j].AddedAt {
			return rows[i].DeviceID < rows[j].DeviceID
		}
		return rows[i].AddedAt > rows[j].AddedAt
	})

	devices := make([]domain.Device, 0, len(rows))
	for _, row := range rows {
		devices = append(devices, c.StorageDeviceToDomain(row))
	}
	return devices
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

func (c *ConvToDamain) StorageThresholdsToDomain(rows []storage.DeviceThreshold) []domain.DeviceThreshold {
	thresholds := make([]domain.DeviceThreshold, 0, len(rows))
	for _, row := range rows {
		thresholds = append(thresholds, c.StorageThresholdToDomain(row))
	}
	return thresholds
}
