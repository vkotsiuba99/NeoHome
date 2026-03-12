package alert

import (
	"context"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) ResolveAlert(ctx context.Context, cmd domain.AlertResolve) (domain.Alert, error) {
	if cmd.UserID <= 0 || cmd.AlertID <= 0 {
		return domain.Alert{}, service.ErrValidation
	}

	Alert, err := svc.alertRepo.GetAlert(ctx, cmd.AlertID)
	if err != nil {
		return domain.Alert{}, mapStorageErr(err)
	}

	alert := svc.toDomain.StorageAlertToDomain(Alert)
	Device, err := svc.deviceRepo.GetDevice(ctx, alert.DeviceID)
	if err != nil {
		return domain.Alert{}, mapStorageErr(err)
	}
	if Device.UserID != cmd.UserID {
		return domain.Alert{}, service.ErrForbidden
	}
	if alert.IsResolved {
		return alert, nil
	}

	alert.IsResolved = true
	alert.ResolvedAt = time.Now().UTC().UnixMilli()
	if err := svc.alertRepo.UpdateAlert(ctx, svc.toStorage.DomainAlertToStorage(alert)); err != nil {
		return domain.Alert{}, mapStorageErr(err)
	}

	return alert, nil
}
