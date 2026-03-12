package alert

import (
	"context"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type DeviceRepo interface {
	ListDevicesByUser(ctx context.Context, userID int64) ([]storage.Device, error)
	GetDevice(ctx context.Context, deviceID int64) (storage.Device, error)
}

//revive:disable-next-line:exported
type AlertRepo interface {
	GetAlert(ctx context.Context, alertID int64) (storage.Alert, error)
	UpdateAlert(ctx context.Context, alert storage.Alert) error
	ListAlerts(ctx context.Context, locationID int64, fromTimestamp int64, hasFrom bool, toTimestamp int64, hasTo bool) ([]storage.Alert, error)
}

type Service struct {
	deviceRepo DeviceRepo
	alertRepo  AlertRepo
	toDomain   ConvToDamain
	toStorage  ConvToStore
	log        slog.Logger
}

func New(deviceRepo DeviceRepo, alertRepo AlertRepo, log *slog.Logger) *Service {
	return &Service{
		deviceRepo: deviceRepo,
		alertRepo:  alertRepo,
		toDomain:   ConvToDamain{},
		toStorage:  ConvToStore{},
		log:        *log.With("domain", "alert"),
	}
}
