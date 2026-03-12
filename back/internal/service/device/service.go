package device

import (
	"context"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

const defaultDeviceStatus = "offline"

type UserRepo interface {
	GetUserByID(ctx context.Context, userID int64) (storage.User, error)
}

type Repo interface {
	CreateDevice(ctx context.Context, device storage.Device) error
	ListDevicesByUser(ctx context.Context, userID int64) ([]storage.Device, error)
	GetDevice(ctx context.Context, deviceID int64) (storage.Device, error)
	UpdateDevice(ctx context.Context, device storage.Device) error
}

type ThresholdRepo interface {
	PutThresholds(ctx context.Context, deviceID int64, thresholds []storage.DeviceThreshold) error
	GetThresholds(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error)
}

type Service struct {
	userRepo            UserRepo
	deviceRepo          Repo
	deviceThresholdRepo ThresholdRepo
	toDomain            ConvToDamain
	toStorage           ConvToStore
	cfg                 service.Config
	log                 slog.Logger
}

func New(userRepo UserRepo, deviceRepo Repo, deviceThresholdRepo ThresholdRepo, cfg service.Config, log *slog.Logger) *Service {
	return &Service{
		userRepo:            userRepo,
		deviceRepo:          deviceRepo,
		deviceThresholdRepo: deviceThresholdRepo,
		toDomain:            ConvToDamain{},
		toStorage:           ConvToStore{},
		cfg:                 cfg,
		log:                 *log.With("domain", "device"),
	}
}
