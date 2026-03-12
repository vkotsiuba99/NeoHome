package device

import (
	"context"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) UpdateDevice(ctx context.Context, cmd domain.DeviceUpdate) (domain.Device, error) {
	device, err := svc.ownedDevice(ctx, cmd.UserID, cmd.DeviceID)
	if err != nil {
		return domain.Device{}, err
	}

	updated := false
	stringPatches := []struct {
		has      bool
		value    string
		current  *string
		required bool
	}{
		{has: cmd.HasDeviceName, value: cmd.DeviceName, current: &device.DeviceName, required: true},
		{has: cmd.HasDeviceType, value: cmd.DeviceType, current: &device.DeviceType, required: true},
		{has: cmd.HasRoomName, value: cmd.RoomName, current: &device.RoomName},
		{has: cmd.HasLocationName, value: cmd.LocationName, current: &device.LocationName},
		{has: cmd.HasStatus, value: cmd.Status, current: &device.Status, required: true},
	}

	for _, patch := range stringPatches {
		if !patch.has {
			continue
		}

		value := strings.TrimSpace(patch.value)
		if patch.required && len(value) == 0 {
			return domain.Device{}, service.ErrValidation
		}
		if value == *patch.current {
			continue
		}

		*patch.current = value
		updated = true
	}

	if cmd.HasLocationID {
		if cmd.LocationID <= 0 {
			return domain.Device{}, service.ErrValidation
		}
		if cmd.LocationID != device.LocationID {
			device.LocationID = cmd.LocationID
			updated = true
		}
	}
	if !updated {
		return domain.Device{}, service.ErrValidation
	}

	device.UpdatedAt = time.Now().UTC().UnixMilli()
	if err := svc.deviceRepo.UpdateDevice(ctx, svc.toStorage.DomainDeviceToStorage(device)); err != nil {
		return domain.Device{}, mapStorageErr(err)
	}

	return device, nil
}
