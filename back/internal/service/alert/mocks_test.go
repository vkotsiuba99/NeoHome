package alert

import (
	"context"
	"io"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type mockDeviceRepo struct {
	listDevicesFn func(ctx context.Context, userID int64) ([]storage.Device, error)
	getDeviceFn   func(ctx context.Context, deviceID int64) (storage.Device, error)
}

func (m *mockDeviceRepo) ListDevicesByUser(ctx context.Context, userID int64) ([]storage.Device, error) {
	if m.listDevicesFn != nil {
		return m.listDevicesFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockDeviceRepo) GetDevice(ctx context.Context, deviceID int64) (storage.Device, error) {
	if m.getDeviceFn != nil {
		return m.getDeviceFn(ctx, deviceID)
	}
	return storage.Device{}, nil
}

type mockAlertRepo struct {
	getAlertFn    func(ctx context.Context, alertID int64) (storage.Alert, error)
	updateAlertFn func(ctx context.Context, alert storage.Alert) error
	listAlertsFn  func(ctx context.Context, locationID int64, fromTimestamp int64, hasFrom bool, toTimestamp int64, hasTo bool) ([]storage.Alert, error)
	lastUpdated   storage.Alert
}

func (m *mockAlertRepo) GetAlert(ctx context.Context, alertID int64) (storage.Alert, error) {
	if m.getAlertFn != nil {
		return m.getAlertFn(ctx, alertID)
	}
	return storage.Alert{}, nil
}

func (m *mockAlertRepo) UpdateAlert(ctx context.Context, alert storage.Alert) error {
	m.lastUpdated = alert
	if m.updateAlertFn != nil {
		return m.updateAlertFn(ctx, alert)
	}
	return nil
}

func (m *mockAlertRepo) ListAlerts(ctx context.Context, locationID int64, fromTimestamp int64, hasFrom bool, toTimestamp int64, hasTo bool) ([]storage.Alert, error) {
	if m.listAlertsFn != nil {
		return m.listAlertsFn(ctx, locationID, fromTimestamp, hasFrom, toTimestamp, hasTo)
	}
	return nil, nil
}

func newAlertServiceForTest(deviceRepo *mockDeviceRepo, alertRepo *mockAlertRepo) *Service {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(deviceRepo, alertRepo, log)
}
