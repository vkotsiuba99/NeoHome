package device

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) ListDevices(ctx context.Context, userID int64) ([]domain.Device, error) {
	svc.log.Info("list devices started", "method", "ListDevices", "user_id", userID)
	if userID <= 0 {
		return nil, service.ErrValidation
	}
	if _, err := svc.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, mapStorageErr(err)
	}
	deviceDTOs, err := svc.deviceRepo.ListDevicesByUser(ctx, userID)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	return svc.toDomain.StorageDevicesToDomain(deviceDTOs), nil
}
