package alert

import (
	"sort"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToDamain struct{}

func NewDomainConv() *ConvToDamain {
	return &ConvToDamain{}
}

func (c *ConvToDamain) StorageAlertToDomain(dto storage.Alert) domain.Alert {
	return domain.Alert{
		AlertID:    dto.AlertID,
		LocationID: dto.LocationID,
		DeviceID:   dto.DeviceID,
		CreatedAt:  dto.CreatedAt,
		Severity:   dto.Severity,
		Message:    dto.Message,
		IsResolved: dto.IsResolved,
		ResolvedAt: dto.ResolvedAt,
	}
}

func (c *ConvToDamain) StorageAlertsToDomain(rows []storage.Alert, allowedDeviceIDs map[int64]struct{}) []domain.Alert {
	alerts := make([]domain.Alert, 0, len(rows))
	for _, dto := range rows {
		if _, ok := allowedDeviceIDs[dto.DeviceID]; !ok {
			continue
		}
		alerts = append(alerts, c.StorageAlertToDomain(dto))
	}
	sort.Slice(alerts, func(i, j int) bool {
		if alerts[i].CreatedAt == alerts[j].CreatedAt {
			return alerts[i].AlertID < alerts[j].AlertID
		}
		return alerts[i].CreatedAt > alerts[j].CreatedAt
	})
	return alerts
}

func (c *ConvToDamain) DevicesToAllowed(devices []storage.Device, locationID int64) (map[int64]struct{}, map[int64]struct{}, bool) {
	allowedDeviceIDs := make(map[int64]struct{}, len(devices))
	allowedLocations := make(map[int64]struct{}, len(devices))
	for _, dto := range devices {
		allowedDeviceIDs[dto.DeviceID] = struct{}{}
		if dto.LocationID > 0 {
			allowedLocations[dto.LocationID] = struct{}{}
		}
	}

	if locationID > 0 {
		_, ok := allowedLocations[locationID]
		return allowedDeviceIDs, allowedLocations, ok
	}

	return allowedDeviceIDs, allowedLocations, true
}
