package telemetry

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

func (svc *Service) makeAlerts(ctx context.Context, device domain.Device, telemetry domain.Telemetry) ([]domain.Alert, error) {
	rows, err := svc.deviceThresholdRepo.GetThresholds(ctx, telemetry.DeviceID)
	if err != nil {
		return nil, mapStorageErr(err)
	}

	alerts := make([]domain.Alert, 0, 1)
	for _, dto := range rows {
		threshold := svc.toDomain.StorageThresholdToDomain(dto)
		metricType := strings.ToLower(strings.TrimSpace(threshold.MetricType))
		if metricType != telemetry.MetricType {
			continue
		}
		switch {
		case threshold.HasMinValue && telemetry.MetricValue < threshold.MinValue:
		case threshold.HasMaxValue && telemetry.MetricValue > threshold.MaxValue:
		default:
			continue
		}

		alertID := int64(0)
		for alertID <= 0 {
			rawID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
			if err != nil {
				return nil, err
			}
			alertID = rawID.Int64()
		}

		severity := strings.ToLower(strings.TrimSpace(dto.Severity))
		if len(severity) == 0 {
			severity = svc.cfg.DefaultSeverity
		}

		message := svc.alertText(telemetry, threshold)
		alert := svc.toDomain.TelemetryToDomainAlert(alertID, device, telemetry, severity, time.Now().UTC().UnixMilli(), message)
		if err := svc.alertRepo.CreateAlert(ctx, svc.toStorage.DomainAlertToStorage(alert)); err != nil {
			return nil, mapStorageErr(err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (svc *Service) alertText(telemetry domain.Telemetry, threshold domain.DeviceThreshold) string {
	bounds := make([]string, 0, 2)
	if threshold.HasMinValue {
		bounds = append(bounds, fmt.Sprintf("min=%.2f", threshold.MinValue))
	}
	if threshold.HasMaxValue {
		bounds = append(bounds, fmt.Sprintf("max=%.2f", threshold.MaxValue))
	}

	boundary := strings.Join(bounds, ", ")
	if len(boundary) == 0 {
		boundary = "n/a"
	}

	unit := strings.TrimSpace(telemetry.Unit)
	if len(unit) > 0 {
		unit = " " + unit
	}

	return fmt.Sprintf("%s is out of threshold: value=%.2f%s (%s)", telemetry.MetricType, telemetry.MetricValue, unit, boundary)
}
