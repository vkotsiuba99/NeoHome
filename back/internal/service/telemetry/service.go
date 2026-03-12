package telemetry

import (
	"context"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

const activeDeviceStatus = "online"

type DeviceRepo interface {
	GetDevice(ctx context.Context, deviceID int64) (storage.Device, error)
	UpdateDevice(ctx context.Context, device storage.Device) error
}

type ThresholdRepo interface {
	GetThresholds(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error)
}

type Repo interface {
	AddTelemetry(ctx context.Context, telemetry storage.Telemetry) error
	ListTelemetry(ctx context.Context, deviceID int64, metricType string, fromTimestamp int64, hasFrom bool, toTimestamp int64, hasTo bool, limit int64) ([]storage.Telemetry, error)
	GetLatestTelemetry(ctx context.Context, deviceID int64) ([]storage.Telemetry, error)
}

type AlertRepo interface {
	CreateAlert(ctx context.Context, alert storage.Alert) error
}

type Service struct {
	deviceRepo          DeviceRepo
	deviceThresholdRepo ThresholdRepo
	telemetryRepo       Repo
	alertRepo           AlertRepo
	toDomain            ConvToDamain
	toStorage           ConvToStore
	cfg                 service.Config
	log                 slog.Logger
}

func New(deviceRepo DeviceRepo, deviceThresholdRepo ThresholdRepo, telemetryRepo Repo, alertRepo AlertRepo, cfg service.Config, log *slog.Logger) *Service {
	return &Service{
		deviceRepo:          deviceRepo,
		deviceThresholdRepo: deviceThresholdRepo,
		telemetryRepo:       telemetryRepo,
		alertRepo:           alertRepo,
		toDomain:            ConvToDamain{},
		toStorage:           ConvToStore{},
		cfg:                 cfg,
		log:                 *log.With("domain", "telemetry"),
	}
}
