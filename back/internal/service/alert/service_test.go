package alert

import (
	"context"
	"errors"
	"testing"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func TestAlertConvertersAndNew(t *testing.T) {
	if NewDomainConv() == nil {
		t.Fatal("NewDomainConv returned nil")
	}
	if NewStoreConv() == nil {
		t.Fatal("NewStoreConv returned nil")
	}

	svc := newAlertServiceForTest(&mockDeviceRepo{}, &mockAlertRepo{})
	if svc == nil {
		t.Fatal("New returned nil")
	}

	dto := storage.Alert{
		AlertID:    1,
		LocationID: 2,
		DeviceID:   3,
		CreatedAt:  4,
		Severity:   "critical",
		Message:    "x",
		IsResolved: true,
		ResolvedAt: 5,
	}
	alert := svc.toDomain.StorageAlertToDomain(dto)
	if alert.AlertID != dto.AlertID || alert.DeviceID != dto.DeviceID || !alert.IsResolved {
		t.Fatalf("unexpected alert conversion: %+v", alert)
	}
	if svc.toStorage.DomainAlertToStorage(alert).AlertID != dto.AlertID {
		t.Fatal("unexpected storage conversion")
	}
}

func TestListAlerts(t *testing.T) {
	svc := newAlertServiceForTest(&mockDeviceRepo{}, &mockAlertRepo{})
	if _, err := svc.ListAlerts(context.Background(), domain.AlertQuery{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for empty user, got %v", err)
	}
	if _, err := svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1, HasFromCreatedAt: true, FromCreatedAt: 10, HasToCreatedAt: true, ToCreatedAt: 1}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for wrong range, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newAlertServiceForTest(&mockDeviceRepo{
			listDevicesFn: func(context.Context, int64) ([]storage.Device, error) { return nil, tc.err },
		}, &mockAlertRepo{})
		_, err := svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1})
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

	svc = newAlertServiceForTest(&mockDeviceRepo{
		listDevicesFn: func(context.Context, int64) ([]storage.Device, error) {
			return []storage.Device{{DeviceID: 1, LocationID: 10}}, nil
		},
	}, &mockAlertRepo{})
	items, err := svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1, LocationID: 999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %+v", items)
	}

	svc = newAlertServiceForTest(&mockDeviceRepo{
		listDevicesFn: func(context.Context, int64) ([]storage.Device, error) {
			return []storage.Device{{DeviceID: 1, LocationID: 10}}, nil
		},
	}, &mockAlertRepo{
		listAlertsFn: func(context.Context, int64, int64, bool, int64, bool) ([]storage.Alert, error) {
			return []storage.Alert{
				{AlertID: 2, DeviceID: 1, CreatedAt: 100},
				{AlertID: 1, DeviceID: 1, CreatedAt: 100},
			}, nil
		},
	})
	items, err = svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1, LocationID: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0].AlertID != 1 || items[1].AlertID != 2 {
		t.Fatalf("unexpected location-scoped result: %+v", items)
	}

	svc = newAlertServiceForTest(&mockDeviceRepo{
		listDevicesFn: func(context.Context, int64) ([]storage.Device, error) { return []storage.Device{}, nil },
	}, &mockAlertRepo{})
	items, err = svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list for no locations, got %+v", items)
	}

	svc = newAlertServiceForTest(&mockDeviceRepo{
		listDevicesFn: func(context.Context, int64) ([]storage.Device, error) {
			return []storage.Device{
				{DeviceID: 1, LocationID: 10},
				{DeviceID: 2, LocationID: 20},
			}, nil
		},
	}, &mockAlertRepo{
		listAlertsFn: func(_ context.Context, locationID int64, _ int64, _ bool, _ int64, _ bool) ([]storage.Alert, error) {
			if locationID == 20 {
				return nil, storage.ErrConflict
			}
			return []storage.Alert{{AlertID: 1, DeviceID: 1, CreatedAt: 1}}, nil
		},
	})
	if _, err := svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected conflict from per-location loop, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newAlertServiceForTest(&mockDeviceRepo{
			listDevicesFn: func(context.Context, int64) ([]storage.Device, error) {
				return []storage.Device{{DeviceID: 1, LocationID: 10}}, nil
			},
		}, &mockAlertRepo{
			listAlertsFn: func(context.Context, int64, int64, bool, int64, bool) ([]storage.Alert, error) {
				return nil, tc.err
			},
		})
		_, err := svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1, LocationID: 10})
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

	svc = newAlertServiceForTest(&mockDeviceRepo{
		listDevicesFn: func(context.Context, int64) ([]storage.Device, error) {
			return []storage.Device{
				{DeviceID: 1, LocationID: 10},
				{DeviceID: 2, LocationID: 20},
			}, nil
		},
	}, &mockAlertRepo{
		listAlertsFn: func(_ context.Context, locationID int64, _ int64, _ bool, _ int64, _ bool) ([]storage.Alert, error) {
			switch locationID {
			case 10:
				return []storage.Alert{
					{AlertID: 3, DeviceID: 1, CreatedAt: 100},
					{AlertID: 2, DeviceID: 999, CreatedAt: 200},
				}, nil
			case 20:
				return []storage.Alert{
					{AlertID: 1, DeviceID: 2, CreatedAt: 100},
					{AlertID: 4, DeviceID: 2, CreatedAt: 300},
				}, nil
			default:
				return nil, nil
			}
		},
	})
	items, err = svc.ListAlerts(context.Background(), domain.AlertQuery{UserID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected three alerts after filtering, got %+v", items)
	}
	if items[0].AlertID != 4 || items[1].AlertID != 1 || items[2].AlertID != 3 {
		t.Fatalf("unexpected sorting/filter result: %+v", items)
	}
}

func TestResolveAlert(t *testing.T) {
	svc := newAlertServiceForTest(&mockDeviceRepo{}, &mockAlertRepo{})
	if _, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{}); !errors.Is(err, service.ErrValidation) {
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
		svc = newAlertServiceForTest(&mockDeviceRepo{}, &mockAlertRepo{
			getAlertFn: func(context.Context, int64) (storage.Alert, error) { return storage.Alert{}, tc.err },
		})
		_, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2})
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

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newAlertServiceForTest(&mockDeviceRepo{
			getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{}, tc.err },
		}, &mockAlertRepo{
			getAlertFn: func(context.Context, int64) (storage.Alert, error) {
				return storage.Alert{AlertID: 2, DeviceID: 9}, nil
			},
		})
		_, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2})
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

	svc = newAlertServiceForTest(&mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{UserID: 2}, nil },
	}, &mockAlertRepo{
		getAlertFn: func(context.Context, int64) (storage.Alert, error) {
			return storage.Alert{AlertID: 2, DeviceID: 9}, nil
		},
	})
	if _, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2}); !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}

	svc = newAlertServiceForTest(&mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{UserID: 1}, nil },
	}, &mockAlertRepo{
		getAlertFn: func(context.Context, int64) (storage.Alert, error) {
			return storage.Alert{AlertID: 2, DeviceID: 9, IsResolved: true, ResolvedAt: 7}, nil
		},
	})
	alert, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2})
	if err != nil || !alert.IsResolved || alert.ResolvedAt != 7 {
		t.Fatalf("unexpected already-resolved result: %+v err=%v", alert, err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newAlertServiceForTest(&mockDeviceRepo{
			getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{UserID: 1}, nil },
		}, &mockAlertRepo{
			getAlertFn: func(context.Context, int64) (storage.Alert, error) {
				return storage.Alert{AlertID: 2, DeviceID: 9}, nil
			},
			updateAlertFn: func(context.Context, storage.Alert) error { return tc.err },
		})
		_, err := svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2})
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

	alertRepo := &mockAlertRepo{
		getAlertFn: func(context.Context, int64) (storage.Alert, error) {
			return storage.Alert{AlertID: 2, DeviceID: 9}, nil
		},
	}
	svc = newAlertServiceForTest(&mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{UserID: 1}, nil },
	}, alertRepo)
	alert, err = svc.ResolveAlert(context.Background(), domain.AlertResolve{UserID: 1, AlertID: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !alert.IsResolved || alert.ResolvedAt <= 0 || !alertRepo.lastUpdated.IsResolved {
		t.Fatalf("unexpected resolve result: %+v, updated=%+v", alert, alertRepo.lastUpdated)
	}
}
