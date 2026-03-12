package telemetry

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type mockDeviceRepo struct {
	getDeviceFn      func(ctx context.Context, deviceID int64) (storage.Device, error)
	updateDeviceFn   func(ctx context.Context, device storage.Device) error
	lastUpdateDevice storage.Device
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
	getThresholdsFn func(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error)
}

func (m *mockThresholdRepo) GetThresholds(ctx context.Context, deviceID int64) ([]storage.DeviceThreshold, error) {
	if m.getThresholdsFn != nil {
		return m.getThresholdsFn(ctx, deviceID)
	}
	return nil, nil
}

type mockTelemetryRepo struct {
	addTelemetryFn     func(ctx context.Context, telemetry storage.Telemetry) error
	listTelemetryFn    func(ctx context.Context, deviceID int64, metricType string, fromTimestamp int64, hasFrom bool, toTimestamp int64, hasTo bool, limit int64) ([]storage.Telemetry, error)
	getLatestFn        func(ctx context.Context, deviceID int64) ([]storage.Telemetry, error)
	lastAddTelemetry   storage.Telemetry
	lastListMetricType string
	lastListLimit      int64
}

func (m *mockTelemetryRepo) AddTelemetry(ctx context.Context, telemetry storage.Telemetry) error {
	m.lastAddTelemetry = telemetry
	if m.addTelemetryFn != nil {
		return m.addTelemetryFn(ctx, telemetry)
	}
	return nil
}

func (m *mockTelemetryRepo) ListTelemetry(
	ctx context.Context,
	deviceID int64,
	metricType string,
	fromTimestamp int64,
	hasFrom bool,
	toTimestamp int64,
	hasTo bool,
	limit int64,
) ([]storage.Telemetry, error) {
	m.lastListMetricType = metricType
	m.lastListLimit = limit
	if m.listTelemetryFn != nil {
		return m.listTelemetryFn(ctx, deviceID, metricType, fromTimestamp, hasFrom, toTimestamp, hasTo, limit)
	}
	return nil, nil
}

func (m *mockTelemetryRepo) GetLatestTelemetry(ctx context.Context, deviceID int64) ([]storage.Telemetry, error) {
	if m.getLatestFn != nil {
		return m.getLatestFn(ctx, deviceID)
	}
	return nil, nil
}

type mockAlertRepo struct {
	createAlertFn   func(ctx context.Context, alert storage.Alert) error
	lastCreateAlert storage.Alert
}

func (m *mockAlertRepo) CreateAlert(ctx context.Context, alert storage.Alert) error {
	m.lastCreateAlert = alert
	if m.createAlertFn != nil {
		return m.createAlertFn(ctx, alert)
	}
	return nil
}

func newTelemetryServiceForTest(
	deviceRepo *mockDeviceRepo,
	thresholdRepo *mockThresholdRepo,
	telemetryRepo *mockTelemetryRepo,
	alertRepo *mockAlertRepo,
) *Service {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(deviceRepo, thresholdRepo, telemetryRepo, alertRepo, service.Config{
		JWTSecret:         []byte("test-secret"),
		TokenTTL:          time.Hour,
		MinPasswordLength: 8,
		DefaultSeverity:   "critical",
		MaxHistoryLimit:   1000,
		DefaultHistory:    100,
	}, log)
}
