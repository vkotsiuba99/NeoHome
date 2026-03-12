package handler

import (
	"context"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

const (
	defaultTelemetryLimit = 100
	maxTelemetryLimit     = 1000
)

type AuthService interface {
	CreateUser(ctx context.Context, cmd domain.Register) (domain.User, error)
	LoginEmail(ctx context.Context, cmd domain.Auth) (domain.Session, error)
	GetUser(ctx context.Context, userID int64) (domain.User, error)
	UpdateUser(ctx context.Context, cmd domain.Update) (domain.User, error)
}

type DeviceService interface {
	ListDevices(ctx context.Context, userID int64) ([]domain.Device, error)
	CreateDevice(ctx context.Context, cmd domain.DeviceCreate) (domain.Device, error)
	UpdateDevice(ctx context.Context, cmd domain.DeviceUpdate) (domain.Device, error)
	GetThresholds(ctx context.Context, userID int64, deviceID int64) ([]domain.DeviceThreshold, error)
	PutThresholds(ctx context.Context, cmd domain.ThresholdsUpsert) ([]domain.DeviceThreshold, error)
}

type TelemetryService interface {
	IngestTelemetry(ctx context.Context, cmd domain.TelemetryIngest) (domain.IngestResult, error)
	ListTelemetry(ctx context.Context, query domain.TelemetryQuery) ([]domain.Telemetry, error)
	GetLatest(ctx context.Context, userID int64, deviceID int64) ([]domain.Telemetry, error)
}

type AlertService interface {
	ListAlerts(ctx context.Context, query domain.AlertQuery) ([]domain.Alert, error)
	ResolveAlert(ctx context.Context, cmd domain.AlertResolve) (domain.Alert, error)
}

type Handler struct {
	authService      AuthService
	deviceService    DeviceService
	telemetryService TelemetryService
	alertService     AlertService
	ReqConv          *ReqConv
	RespConv         *RespConv
	log              slog.Logger
}

func New(authService AuthService, deviceService DeviceService, telemetryService TelemetryService, alertService AlertService, mqttTopicTemplate string, log *slog.Logger) *Handler {
	return &Handler{
		authService:      authService,
		deviceService:    deviceService,
		telemetryService: telemetryService,
		alertService:     alertService,
		ReqConv:          NewReqConv(mqttTopicTemplate),
		RespConv:         NewRespConv(),
		log: *log.With(
			"layer", "transport",
			"component", "handler",
		),
	}
}
