package telemetry

import (
	"context"
	"strings"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) ListTelemetry(ctx context.Context, query domain.TelemetryQuery) ([]domain.Telemetry, error) {
	if _, err := svc.ownedDevice(ctx, query.UserID, query.DeviceID); err != nil {
		return nil, err
	}
	if query.HasFromRecordedAt && query.HasToRecordedAt && query.FromRecordedAt > query.ToRecordedAt {
		return nil, service.ErrValidation
	}

	limit := query.Limit
	switch {
	case limit <= 0:
		limit = svc.cfg.DefaultHistory
	case limit > svc.cfg.MaxHistoryLimit:
		limit = svc.cfg.MaxHistoryLimit
	}

	rows, err := svc.telemetryRepo.ListTelemetry(ctx, query.DeviceID, strings.ToLower(strings.TrimSpace(query.MetricType)), query.FromRecordedAt, query.HasFromRecordedAt, query.ToRecordedAt, query.HasToRecordedAt, limit)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	return svc.toDomain.StorageTelemetryListToDomain(rows), nil
}
