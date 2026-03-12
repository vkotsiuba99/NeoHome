package device

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type mockUserRepo struct {
	getUserByIDFn func(ctx context.Context, userID int64) (storage.User, error)
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, userID int64) (storage.User, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, userID)
	}
	return storage.User{}, nil
}

type mockDeviceRepo struct {
	createDeviceFn      func(ctx context.Context, device storage.Device) error
	listDevicesByUserFn func(ctx context.Context, userID int64) ([]storage.Device, error)
	getDeviceFn         func(ctx context.Context, deviceID int64) (storage.Device, error)
	updateDeviceFn      func(ctx context.Context, device storage.Device) error
	lastCreateDevice    storage.Device
	lastUpdateDevice    storage.Device
}

func (m *mockDeviceRepo) CreateDevice(ctx context.Context, device storage.Device) error {
	m.lastCreateDevice = device
	if m.createDeviceFn != nil {
		return m.createDeviceFn(ctx, device)
	}
	return nil
}

func (m *mockDeviceRepo) ListDevicesByUser(ctx context.Context, userID int64) ([]storage.Device, error) {
	if m.listDevicesByUserFn != nil {
		return m.listDevicesByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockDeviceRepo) GetDevice(ctx context.Context, deviceID int64) (storage.Device, error) {
	if m.getDeviceFn != nil {
		return m.getDeviceFn(ctx, deviceID)
	}
	return storage.Device{}, nil
}

func (m *mockDeviceRepo) UpdateDevice(ctx context.Context, device storage.Device) error {
	m.lastUpdateDevice = device
	if m.updateDeviceFn != nil {
		return m.updateDeviceFn(ctx, device)
	}
	return nil
}

type mockThresholdRepo struct {
	putThresholdsFn   func(ctx context.Context, deviceID int64, thresholds []storage.DeviceThreshold) error
	getThresholdsFn   func(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error)
	lastPutDeviceID   int64
	lastPutThresholds []storage.DeviceThreshold
}

func (m *mockThresholdRepo) PutThresholds(ctx context.Context, deviceID int64, thresholds []storage.DeviceThreshold) error {
	m.lastPutDeviceID = deviceID
	m.lastPutThresholds = thresholds
	if m.putThresholdsFn != nil {
		return m.putThresholdsFn(ctx, deviceID, thresholds)
	}
	return nil
}

func (m *mockThresholdRepo) GetThresholds(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error) {
	if m.getThresholdsFn != nil {
		return m.getThresholdsFn(ctx, deviceID)
	}
	return nil, nil
}

func newDeviceServiceForTest(userRepo *mockUserRepo, deviceRepo *mockDeviceRepo, thresholdRepo *mockThresholdRepo) *Service {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(userRepo, deviceRepo, thresholdRepo, service.Config{
		JWTSecret:         []byte("test-secret"),
		TokenTTL:          time.Hour,
		MinPasswordLength: 8,
		DefaultSeverity:   "critical",
		MaxHistoryLimit:   1000,
		DefaultHistory:    100,
	}, log)
}
