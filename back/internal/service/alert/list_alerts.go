package alert

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func (svc *Service) ListAlerts(ctx context.Context, query domain.AlertQuery) ([]domain.Alert, error) {
	if query.UserID <= 0 {
		return nil, service.ErrValidation
	}
	if query.HasFromCreatedAt && query.HasToCreatedAt && query.FromCreatedAt > query.ToCreatedAt {
		return nil, service.ErrValidation
	}

	devices, err := svc.deviceRepo.ListDevicesByUser(ctx, query.UserID)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	allowedDeviceIDs, allowedLocations, locationAllowed := svc.toDomain.DevicesToAllowed(devices, query.LocationID)
	if !locationAllowed {
		return []domain.Alert{}, nil
	}

	rows := make([]storage.Alert, 0, 32)
	if query.LocationID > 0 {
		locationRows, listErr := svc.alertRepo.ListAlerts(ctx, query.LocationID, query.FromCreatedAt, query.HasFromCreatedAt, query.ToCreatedAt, query.HasToCreatedAt)
		if listErr != nil {
			return nil, mapStorageErr(listErr)
		}
		rows = append(rows, locationRows...)
	} else {
		if len(allowedLocations) == 0 {
			return []domain.Alert{}, nil
		}
		for locationID := range allowedLocations {
			locationRows, listErr := svc.alertRepo.ListAlerts(ctx, locationID, query.FromCreatedAt, query.HasFromCreatedAt, query.ToCreatedAt, query.HasToCreatedAt)
			if listErr != nil {
				return nil, mapStorageErr(listErr)
			}
			rows = append(rows, locationRows...)
		}
	}

	return svc.toDomain.StorageAlertsToDomain(rows, allowedDeviceIDs), nil
}
