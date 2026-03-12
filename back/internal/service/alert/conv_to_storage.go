package alert

import (
	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToStore struct{}

func NewStoreConv() *ConvToStore {
	return &ConvToStore{}
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
