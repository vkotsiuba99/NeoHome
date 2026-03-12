package device

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) CreateDevice(ctx context.Context, cmd domain.DeviceCreate) (domain.Device, error) {
	svc.log.Info("create device started", "method", "CreateDevice", "user_id", cmd.UserID)

	cmd = svc.toDomain.CreateToDomainNormalize(cmd)

	if cmd.UserID <= 0 {
		return domain.Device{}, service.ErrValidation
	}
	if len(cmd.DeviceName) == 0 || len(cmd.DeviceType) == 0 {
		return domain.Device{}, service.ErrValidation
	}
	if len(cmd.Status) == 0 {
		cmd.Status = defaultDeviceStatus
	}
	if _, err := svc.userRepo.GetUserByID(ctx, cmd.UserID); err != nil {
		return domain.Device{}, mapStorageErr(err)
	}

	device, err := svc.toDomain.CreateToDomainDeviceAuto(cmd)
	if err != nil {
		return domain.Device{}, err
	}
	if err := svc.deviceRepo.CreateDevice(ctx, svc.toStorage.DomainDeviceToStorage(device)); err != nil {
		return domain.Device{}, mapStorageErr(err)
	}

	return device, nil
}
