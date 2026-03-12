package telemetry

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func TestTelemetryConvertersAndNew(t *testing.T) {
	if NewDomainConv() == nil {
		t.Fatal("NewDomainConv returned nil")
	}
	if NewStoreConv() == nil {
		t.Fatal("NewStoreConv returned nil")
	}

	svc := newTelemetryServiceForTest(&mockDeviceRepo{}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})
	if svc == nil {
		t.Fatal("New returned nil")
	}

	device := svc.toDomain.StorageDeviceToDomain(storage.Device{DeviceID: 1, UserID: 2, DeviceName: "d"})
	if device.DeviceID != 1 || device.UserID != 2 {
		t.Fatalf("unexpected device conversion: %+v", device)
	}

	telemetry := svc.toDomain.IngestToDomainTelemetry(domain.TelemetryIngest{MetricValue: 42}, device, 11, 12, "temp", "C", "kitchen", "home", 90, 50)
	if telemetry.TelemetryID != 11 || telemetry.RecordedAt != 12 || telemetry.MetricType != "temp" {
		t.Fatalf("unexpected new telemetry: %+v", telemetry)
	}

	dto := svc.toStorage.DomainTelemetryToStorage(telemetry)
	telemetry2 := svc.toDomain.StorageTelemetryToDomain(dto)
	if telemetry2.TelemetryID != telemetry.TelemetryID || telemetry2.DeviceID != telemetry.DeviceID {
		t.Fatalf("unexpected telemetry conversion: %+v", telemetry2)
	}

	list := svc.toDomain.StorageTelemetryListToDomain([]storage.Telemetry{{TelemetryID: 1}, {TelemetryID: 2}})
	if len(list) != 2 || list[1].TelemetryID != 2 {
		t.Fatalf("unexpected telemetry list conversion: %+v", list)
	}

	th := svc.toDomain.StorageThresholdToDomain(storage.DeviceThreshold{
		MetricType:  "temp",
		MinValue:    math.NaN(),
		HasMinValue: false,
		MaxValue:    30,
		HasMaxValue: false,
	})
	if th.HasMinValue {
		t.Fatalf("expected min absent, got %+v", th)
	}
	if !th.HasMaxValue {
		t.Fatalf("expected max present via non-NaN, got %+v", th)
	}

	alert := svc.toDomain.TelemetryToDomainAlert(7, domain.Device{LocationID: 99}, telemetry, "critical", 123, "msg")
	if alert.AlertID != 7 || alert.LocationID != 99 || alert.DeviceID != telemetry.DeviceID {
		t.Fatalf("unexpected alert: %+v", alert)
	}
	if svc.toStorage.DomainAlertToStorage(alert).AlertID != 7 {
		t.Fatal("alert storage conversion failed")
	}
}

func TestTelemetryOwnedDevice(t *testing.T) {
	svc := newTelemetryServiceForTest(&mockDeviceRepo{}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})
	if _, err := svc.ownedDevice(context.Background(), 0, 1); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newTelemetryServiceForTest(&mockDeviceRepo{
			getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{}, tc.err },
		}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})
		_, err := svc.ownedDevice(context.Background(), 7, 11)
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	svc = newTelemetryServiceForTest(&mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) {
			return storage.Device{DeviceID: 11, UserID: 8}, nil
		},
	}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})
	if _, err := svc.ownedDevice(context.Background(), 7, 11); !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}

	svc = newTelemetryServiceForTest(&mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) {
			return storage.Device{DeviceID: 11, UserID: 7}, nil
		},
	}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})
	device, err := svc.ownedDevice(context.Background(), 7, 11)
	if err != nil || device.DeviceID != 11 {
		t.Fatalf("unexpected owned device result: %+v, err=%v", device, err)
	}
}

func TestListTelemetry(t *testing.T) {
	deviceRepo := &mockDeviceRepo{getDeviceFn: func(context.Context, int64) (storage.Device, error) {
		return storage.Device{DeviceID: 5, UserID: 1}, nil
	}}
	telemetryRepo := &mockTelemetryRepo{}
	svc := newTelemetryServiceForTest(deviceRepo, &mockThresholdRepo{}, telemetryRepo, &mockAlertRepo{})

	if _, err := svc.ListTelemetry(context.Background(), domain.TelemetryQuery{UserID: 1, DeviceID: 5, HasFromRecordedAt: true, FromRecordedAt: 10, HasToRecordedAt: true, ToRecordedAt: 1}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for invalid time range, got %v", err)
	}
	if _, err := svc.ListTelemetry(context.Background(), domain.TelemetryQuery{UserID: 2, DeviceID: 5, Limit: 1}); !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("expected forbidden from ownedDevice, got %v", err)
	}

	telemetryRepo.listTelemetryFn = func(context.Context, int64, string, int64, bool, int64, bool, int64) ([]storage.Telemetry, error) {
		return []storage.Telemetry{{TelemetryID: 1, MetricType: "temp"}}, nil
	}
	items, err := svc.ListTelemetry(context.Background(), domain.TelemetryQuery{
		UserID:     1,
		DeviceID:   5,
		MetricType: " TEMP ",
		Limit:      0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || telemetryRepo.lastListMetricType != "temp" || telemetryRepo.lastListLimit != svc.cfg.DefaultHistory {
		t.Fatalf("unexpected list call: items=%+v metric=%q limit=%d", items, telemetryRepo.lastListMetricType, telemetryRepo.lastListLimit)
	}

	_, err = svc.ListTelemetry(context.Background(), domain.TelemetryQuery{
		UserID:   1,
		DeviceID: 5,
		Limit:    svc.cfg.MaxHistoryLimit + 1,
	})
	if err != nil {
		t.Fatalf("unexpected error for high limit: %v", err)
	}
	if telemetryRepo.lastListLimit != svc.cfg.MaxHistoryLimit {
		t.Fatalf("expected max history limit clamp, got %d", telemetryRepo.lastListLimit)
	}

	_, err = svc.ListTelemetry(context.Background(), domain.TelemetryQuery{
		UserID:   1,
		DeviceID: 5,
		Limit:    7,
	})
	if err != nil {
		t.Fatalf("unexpected error for explicit limit: %v", err)
	}
	if telemetryRepo.lastListLimit != 7 {
		t.Fatalf("expected explicit limit passthrough, got %d", telemetryRepo.lastListLimit)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		telemetryRepo.listTelemetryFn = func(context.Context, int64, string, int64, bool, int64, bool, int64) ([]storage.Telemetry, error) {
			return nil, tc.err
		}
		_, err := svc.ListTelemetry(context.Background(), domain.TelemetryQuery{UserID: 1, DeviceID: 5, Limit: 1})
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}
}

func TestGetLatest(t *testing.T) {
	deviceRepo := &mockDeviceRepo{getDeviceFn: func(context.Context, int64) (storage.Device, error) {
		return storage.Device{DeviceID: 5, UserID: 1}, nil
	}}
	telemetryRepo := &mockTelemetryRepo{}
	svc := newTelemetryServiceForTest(deviceRepo, &mockThresholdRepo{}, telemetryRepo, &mockAlertRepo{})

	if _, err := svc.GetLatest(context.Background(), 2, 5); !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("expected forbidden from ownedDevice, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		telemetryRepo.getLatestFn = func(context.Context, int64) ([]storage.Telemetry, error) { return nil, tc.err }
		_, err := svc.GetLatest(context.Background(), 1, 5)
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	telemetryRepo.getLatestFn = func(context.Context, int64) ([]storage.Telemetry, error) {
		return []storage.Telemetry{
			{MetricType: "temp", RecordedAt: 9, TelemetryID: 1},
			{MetricType: "temp", RecordedAt: 11, TelemetryID: 1},
			{MetricType: " Temp ", RecordedAt: 10, TelemetryID: 1},
			{MetricType: "temp", RecordedAt: 10, TelemetryID: 2},
			{MetricType: "co2", RecordedAt: 5, TelemetryID: 1},
			{MetricType: "co2", RecordedAt: 5, TelemetryID: 2},
			{MetricType: "humidity", RecordedAt: 20, TelemetryID: 1},
		}, nil
	}
	items, err := svc.GetLatest(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 || items[0].MetricType != "co2" || items[0].TelemetryID != 2 || items[1].MetricType != "humidity" || items[2].TelemetryID != 1 {
		t.Fatalf("unexpected latest telemetry result: %+v", items)
	}
}

func TestAlertText(t *testing.T) {
	svc := newTelemetryServiceForTest(&mockDeviceRepo{}, &mockThresholdRepo{}, &mockTelemetryRepo{}, &mockAlertRepo{})

	msg := svc.alertText(domain.Telemetry{MetricType: "temp", MetricValue: 50, Unit: "C"}, domain.DeviceThreshold{
		HasMinValue: true,
		MinValue:    10,
		HasMaxValue: true,
		MaxValue:    40,
	})
	if !strings.Contains(msg, "min=10.00") || !strings.Contains(msg, "max=40.00") || !strings.Contains(msg, " C") {
		t.Fatalf("unexpected alert text: %q", msg)
	}

	msg = svc.alertText(domain.Telemetry{MetricType: "temp", MetricValue: 50}, domain.DeviceThreshold{})
	if !strings.Contains(msg, "(n/a)") {
		t.Fatalf("unexpected fallback alert text: %q", msg)
	}
}

func TestMakeAlerts(t *testing.T) {
	baseDevice := domain.Device{DeviceID: 5, LocationID: 9}
	baseTelemetry := domain.Telemetry{DeviceID: 5, MetricType: "temp", MetricValue: 50, Unit: "C"}

	thresholdRepo := &mockThresholdRepo{}
	alertRepo := &mockAlertRepo{}
	svc := newTelemetryServiceForTest(&mockDeviceRepo{}, thresholdRepo, &mockTelemetryRepo{}, alertRepo)

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) { return nil, tc.err }
		_, err := svc.makeAlerts(context.Background(), baseDevice, baseTelemetry)
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) {
		return []storage.DeviceThreshold{
			{MetricType: "temp", HasMinValue: true, MinValue: 10, HasMaxValue: true, MaxValue: 100},
		}, nil
	}
	alerts, err := svc.makeAlerts(context.Background(), baseDevice, baseTelemetry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts when value in range, got %+v", alerts)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) {
			return []storage.DeviceThreshold{{MetricType: "temp", HasMaxValue: true, MaxValue: 10}}, nil
		}
		alertRepo.createAlertFn = func(context.Context, storage.Alert) error { return tc.err }
		_, err := svc.makeAlerts(context.Background(), baseDevice, baseTelemetry)
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	alertRepo.createAlertFn = nil
	thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) {
		return []storage.DeviceThreshold{
			{MetricType: " humidity ", HasMaxValue: true, MaxValue: 90},
			{MetricType: "temp", HasMinValue: true, MinValue: 60, Severity: ""},
			{MetricType: "temp", HasMaxValue: true, MaxValue: 40, Severity: " high "},
		}, nil
	}
	alerts, err = svc.makeAlerts(context.Background(), baseDevice, baseTelemetry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 2 {
		t.Fatalf("expected two alerts, got %+v", alerts)
	}
	if alerts[0].Severity != svc.cfg.DefaultSeverity || alerts[1].Severity != "high" {
		t.Fatalf("unexpected severities: %+v", alerts)
	}
}

func TestIngestTelemetry(t *testing.T) {
	baseDevice := storage.Device{
		DeviceID:       5,
		UserID:         1,
		RoomName:       "room",
		LocationName:   "home",
		BatteryLevel:   80,
		SignalStrength: 70,
	}
	deviceRepo := &mockDeviceRepo{}
	thresholdRepo := &mockThresholdRepo{getThresholdsFn: func(context.Context, int64) ([]storage.DeviceThreshold, error) { return nil, nil }}
	telemetryRepo := &mockTelemetryRepo{}
	alertRepo := &mockAlertRepo{}
	svc := newTelemetryServiceForTest(deviceRepo, thresholdRepo, telemetryRepo, alertRepo)

	if _, err := svc.IngestTelemetry(context.Background(), domain.TelemetryIngest{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		deviceRepo.getDeviceFn = func(context.Context, int64) (storage.Device, error) { return storage.Device{}, tc.err }
		_, err := svc.IngestTelemetry(context.Background(), domain.TelemetryIngest{DeviceID: 5, MetricType: "temp"})
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	deviceRepo.getDeviceFn = func(context.Context, int64) (storage.Device, error) { return baseDevice, nil }
	if _, err := svc.IngestTelemetry(context.Background(), domain.TelemetryIngest{
		DeviceID:      5,
		MetricType:    "temp",
		HasRecordedAt: true,
		RecordedAt:    0,
	}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected invalid recordedAt validation, got %v", err)
	}

	for _, tc := range []struct {
		setup    func()
		expected error
		rawErr   string
	}{
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = func(context.Context, storage.Telemetry) error { return storage.ErrNotFound }
			},
			expected: service.ErrNotFound,
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = func(context.Context, storage.Telemetry) error { return storage.ErrConflict }
			},
			expected: service.ErrConflict,
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = func(context.Context, storage.Telemetry) error { return errors.New("add boom") }
			},
			rawErr: "add boom",
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = nil
				deviceRepo.updateDeviceFn = func(context.Context, storage.Device) error { return storage.ErrNotFound }
			},
			expected: service.ErrNotFound,
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = nil
				deviceRepo.updateDeviceFn = func(context.Context, storage.Device) error { return storage.ErrConflict }
			},
			expected: service.ErrConflict,
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = nil
				deviceRepo.updateDeviceFn = func(context.Context, storage.Device) error { return errors.New("update boom") }
			},
			rawErr: "update boom",
		},
		{
			setup: func() {
				telemetryRepo.addTelemetryFn = nil
				deviceRepo.updateDeviceFn = nil
				thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) {
					return nil, errors.New("threshold boom")
				}
			},
			rawErr: "threshold boom",
		},
	} {
		deviceRepo.getDeviceFn = func(context.Context, int64) (storage.Device, error) { return baseDevice, nil }
		deviceRepo.updateDeviceFn = nil
		thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) { return nil, nil }
		tc.setup()
		_, err := svc.IngestTelemetry(context.Background(), domain.TelemetryIngest{DeviceID: 5, MetricType: "temp"})
		if len(tc.rawErr) > 0 {
			if err == nil || err.Error() != tc.rawErr {
				t.Fatalf("expected %q, got %v", tc.rawErr, err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	telemetryRepo.addTelemetryFn = nil
	deviceRepo.updateDeviceFn = nil
	thresholdRepo.getThresholdsFn = func(context.Context, int64) ([]storage.DeviceThreshold, error) {
		return []storage.DeviceThreshold{
			{MetricType: "temp", HasMaxValue: true, MaxValue: 10},
		}, nil
	}
	result, err := svc.IngestTelemetry(context.Background(), domain.TelemetryIngest{
		DeviceID:       5,
		MetricType:     " Temp ",
		MetricValue:    22.5,
		Unit:           " C ",
		RoomName:       " ",
		LocationName:   " ",
		HasRecordedAt:  true,
		RecordedAt:     111,
		HasBattery:     true,
		BatteryLevel:   91,
		HasSignal:      true,
		SignalStrength: 62,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Telemetry.MetricType != "temp" || result.Telemetry.RoomName != baseDevice.RoomName || result.Telemetry.LocationName != baseDevice.LocationName {
		t.Fatalf("unexpected telemetry result: %+v", result.Telemetry)
	}
	if len(result.Alerts) != 1 {
		t.Fatalf("expected one alert, got %+v", result.Alerts)
	}
	if deviceRepo.lastUpdateDevice.Status != activeDeviceStatus || deviceRepo.lastUpdateDevice.BatteryLevel != 91 || deviceRepo.lastUpdateDevice.SignalStrength != 62 {
		t.Fatalf("unexpected device update payload: %+v", deviceRepo.lastUpdateDevice)
	}
}
