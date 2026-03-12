package device

import (
	"context"
	"sort"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

func (svc *Service) GetThresholds(ctx context.Context, userID int64, deviceID int64) ([]domain.DeviceThreshold, error) {
	if _, err := svc.ownedDevice(ctx, userID, deviceID); err != nil {
		return nil, err
	}
	rows, err := svc.deviceThresholdRepo.GetThresholds(ctx, deviceID)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].MetricType < rows[j].MetricType
	})

	return svc.toDomain.StorageThresholdsToDomain(rows), nil
}
