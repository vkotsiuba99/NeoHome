package telemetry

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) ownedDevice(ctx context.Context, userID int64, deviceID int64) (domain.Device, error) {
	if userID <= 0 || deviceID <= 0 {
		return domain.Device{}, service.ErrValidation
	}

	dto, err := svc.deviceRepo.GetDevice(ctx, deviceID)
	if err != nil {
		return domain.Device{}, mapStorageErr(err)
	}

	device := svc.toDomain.StorageDeviceToDomain(dto)
	if device.UserID != userID {
		return domain.Device{}, service.ErrForbidden
	}

	return device, nil
}
