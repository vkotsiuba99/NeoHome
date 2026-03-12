package device

import (
	"context"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) PutThresholds(ctx context.Context, cmd domain.ThresholdsUpsert) ([]domain.DeviceThreshold, error) {
	if _, err := svc.ownedDevice(ctx, cmd.UserID, cmd.DeviceID); err != nil {
		return nil, err
	}
	if len(cmd.Thresholds) == 0 {
		return nil, service.ErrValidation
	}

	now := time.Now().UTC().UnixMilli()
	rows, ok := svc.toStorage.PatchesToStorageThresholds(cmd.DeviceID, cmd.Thresholds, svc.cfg.DefaultSeverity, now)
	if !ok {
		return nil, service.ErrValidation
	}

	if err := svc.deviceThresholdRepo.PutThresholds(ctx, cmd.DeviceID, rows); err != nil {
		return nil, mapStorageErr(err)
	}

	return svc.GetThresholds(ctx, cmd.UserID, cmd.DeviceID)
}
