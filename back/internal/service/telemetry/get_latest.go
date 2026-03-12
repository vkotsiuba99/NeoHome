package telemetry

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

func (svc *Service) GetLatest(ctx context.Context, userID int64, deviceID int64) ([]domain.Telemetry, error) {
	if _, err := svc.ownedDevice(ctx, userID, deviceID); err != nil {
		return nil, err
	}

	rows, err := svc.telemetryRepo.GetLatestTelemetry(ctx, deviceID)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	latestRows := svc.toDomain.StorageTelemetryLatestByMetric(rows)

	return svc.toDomain.StorageTelemetryListToDomain(latestRows), nil
}
