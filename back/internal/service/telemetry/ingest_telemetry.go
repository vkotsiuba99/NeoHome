package telemetry

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) IngestTelemetry(ctx context.Context, cmd domain.TelemetryIngest) (domain.IngestResult, error) {
	normalized, ok := svc.toDomain.IngestToDomainNormalize(cmd)
	if !ok {
		return domain.IngestResult{}, service.ErrValidation
	}

	dto, err := svc.deviceRepo.GetDevice(ctx, normalized.DeviceID)
	if err != nil {
		return domain.IngestResult{}, mapStorageErr(err)
	}

	device := svc.toDomain.StorageDeviceToDomain(dto)
	telemetry, updatedDevice, err := svc.toDomain.IngestToDomainTelemetryAndDevice(normalized, device, activeDeviceStatus)
	if err != nil {
		return domain.IngestResult{}, err
	}

	if err := svc.telemetryRepo.AddTelemetry(ctx, svc.toStorage.DomainTelemetryToStorage(telemetry)); err != nil {
		return domain.IngestResult{}, mapStorageErr(err)
	}

	if err := svc.deviceRepo.UpdateDevice(ctx, svc.toStorage.DomainDeviceToStorage(updatedDevice)); err != nil {
		return domain.IngestResult{}, mapStorageErr(err)
	}

	alerts, err := svc.makeAlerts(ctx, updatedDevice, telemetry)
	if err != nil {
		return domain.IngestResult{}, err
	}

	return domain.IngestResult{Telemetry: telemetry, Alerts: alerts}, nil
}
