package device

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func TestDeviceConvertersAndNew(t *testing.T) {
	if NewDomainConv() == nil {
		t.Fatal("NewDomainConv returned nil")
	}
	if NewStoreConv() == nil {
		t.Fatal("NewStoreConv returned nil")
	}

	svc := newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{}, &mockThresholdRepo{})
	if svc == nil {
		t.Fatal("New returned nil")
	}

	device := svc.toDomain.CreateToDomainDevice(domain.DeviceCreate{
		UserID:       7,
		DeviceName:   "a",
		DeviceType:   "sensor",
		RoomName:     "room",
		LocationID:   10,
		LocationName: "home",
		Status:       "on",
	}, 100, 200)
	if device.DeviceID != 100 || device.UserID != 7 || device.AddedAt != 200 || device.UpdatedAt != 200 {
		t.Fatalf("unexpected new device: %+v", device)
	}

	dto := svc.toStorage.DomainDeviceToStorage(device)
	back := svc.toDomain.StorageDeviceToDomain(dto)
	if back.DeviceID != device.DeviceID || back.LocationName != device.LocationName {
		t.Fatalf("unexpected conversion result: %+v", back)
	}

	devs := svc.toDomain.StorageDevicesToDomain([]storage.Device{{DeviceID: 1}, {DeviceID: 2}})
	if len(devs) != 2 || devs[1].DeviceID != 2 {
		t.Fatalf("unexpected devices conversion: %+v", devs)
	}

	th := svc.toStorage.PatchToStorageThreshold(1, domain.ThresholdPatch{
		MetricType:  "temp",
		HasMinValue: false,
		HasMaxValue: true,
		MaxValue:    10,
	}, "temp", "critical", 123)
	if !math.IsNaN(th.MinValue) || th.MaxValue != 10 || !th.HasMaxValue {
		t.Fatalf("unexpected threshold dto: %+v", th)
	}

	domainThreshold := svc.toDomain.StorageThresholdToDomain(storage.DeviceThreshold{
		DeviceID:    1,
		MetricType:  "temp",
		MinValue:    math.NaN(),
		HasMinValue: false,
		MaxValue:    15,
		HasMaxValue: false,
		Severity:    "warn",
		UpdatedAt:   321,
	})
	if domainThreshold.HasMinValue {
		t.Fatalf("expected min value to be absent: %+v", domainThreshold)
	}
	if !domainThreshold.HasMaxValue {
		t.Fatalf("expected max value to be considered present: %+v", domainThreshold)
	}
}

func TestOwnedDevice(t *testing.T) {
	svc := newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{}, &mockThresholdRepo{})
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
		svc = newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{
			getDeviceFn: func(context.Context, int64) (storage.Device, error) { return storage.Device{}, tc.err },
		}, &mockThresholdRepo{})
		_, err := svc.ownedDevice(context.Background(), 7, 10)
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

	svc = newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) {
			return storage.Device{DeviceID: 10, UserID: 8}, nil
		},
	}, &mockThresholdRepo{})
	if _, err := svc.ownedDevice(context.Background(), 7, 10); !errors.Is(err, service.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}

	svc = newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) {
			return storage.Device{DeviceID: 10, UserID: 7}, nil
		},
	}, &mockThresholdRepo{})
	device, err := svc.ownedDevice(context.Background(), 7, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if device.DeviceID != 10 || device.UserID != 7 {
		t.Fatalf("unexpected device: %+v", device)
	}
}

func TestCreateDevice(t *testing.T) {
	baseUserRepo := &mockUserRepo{getUserByIDFn: func(context.Context, int64) (storage.User, error) { return storage.User{UserID: 1}, nil }}
	svc := newDeviceServiceForTest(baseUserRepo, &mockDeviceRepo{}, &mockThresholdRepo{})

	if _, err := svc.CreateDevice(context.Background(), domain.DeviceCreate{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for empty command, got %v", err)
	}
	if _, err := svc.CreateDevice(context.Background(), domain.DeviceCreate{UserID: 1, DeviceName: "a"}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for missing type, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		svc = newDeviceServiceForTest(&mockUserRepo{
			getUserByIDFn: func(context.Context, int64) (storage.User, error) { return storage.User{}, tc.err },
		}, &mockDeviceRepo{}, &mockThresholdRepo{})
		_, err := svc.CreateDevice(context.Background(), domain.DeviceCreate{UserID: 1, DeviceName: "a", DeviceType: "sensor"})
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
		svc = newDeviceServiceForTest(baseUserRepo, &mockDeviceRepo{
			createDeviceFn: func(context.Context, storage.Device) error { return tc.err },
		}, &mockThresholdRepo{})
		_, err := svc.CreateDevice(context.Background(), domain.DeviceCreate{UserID: 1, DeviceName: "a", DeviceType: "sensor"})
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

	deviceRepo := &mockDeviceRepo{}
	svc = newDeviceServiceForTest(baseUserRepo, deviceRepo, &mockThresholdRepo{})
	created, err := svc.CreateDevice(context.Background(), domain.DeviceCreate{
		UserID:       1,
		DeviceName:   "  Sensor ",
		DeviceType:   "  temp ",
		RoomName:     "  kitchen ",
		LocationName: "  home ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Status != defaultDeviceStatus || created.LocationID <= 0 {
		t.Fatalf("unexpected created device: %+v", created)
	}
	if deviceRepo.lastCreateDevice.DeviceName != "Sensor" || deviceRepo.lastCreateDevice.RoomName != "kitchen" {
		t.Fatalf("unexpected DTO normalization: %+v", deviceRepo.lastCreateDevice)
	}
}

func TestUpdateDevice(t *testing.T) {
	base := storage.Device{
		DeviceID:     5,
		UserID:       1,
		DeviceName:   "old",
		DeviceType:   "sensor",
		RoomName:     "room",
		LocationID:   9,
		LocationName: "loc",
		Status:       "offline",
		AddedAt:      1,
		UpdatedAt:    1,
	}

	makeService := func(updateErr error) *Service {
		return newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{
			getDeviceFn:    func(context.Context, int64) (storage.Device, error) { return base, nil },
			updateDeviceFn: func(context.Context, storage.Device) error { return updateErr },
		}, &mockThresholdRepo{})
	}

	svc := makeService(nil)
	if _, err := svc.UpdateDevice(context.Background(), domain.DeviceUpdate{UserID: 0, DeviceID: 5}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation from ownedDevice, got %v", err)
	}

	for _, cmd := range []domain.DeviceUpdate{
		{UserID: 1, DeviceID: 5, HasDeviceName: true, DeviceName: " "},
		{UserID: 1, DeviceID: 5, HasDeviceType: true, DeviceType: " "},
		{UserID: 1, DeviceID: 5, HasLocationID: true, LocationID: 0},
		{UserID: 1, DeviceID: 5, HasStatus: true, Status: " "},
		{UserID: 1, DeviceID: 5, HasRoomName: true, RoomName: "room"},
	} {
		if _, err := svc.UpdateDevice(context.Background(), cmd); !errors.Is(err, service.ErrValidation) {
			t.Fatalf("expected validation, got %v for cmd %+v", err, cmd)
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
		svc = makeService(tc.err)
		_, err := svc.UpdateDevice(context.Background(), domain.DeviceUpdate{
			UserID:        1,
			DeviceID:      5,
			HasDeviceName: true,
			DeviceName:    "new",
		})
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

	deviceRepo := &mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) { return base, nil },
	}
	svc = newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, &mockThresholdRepo{})
	updated, err := svc.UpdateDevice(context.Background(), domain.DeviceUpdate{
		UserID:          1,
		DeviceID:        5,
		HasDeviceName:   true,
		DeviceName:      " new ",
		HasDeviceType:   true,
		DeviceType:      " meter ",
		HasRoomName:     true,
		RoomName:        " hall ",
		HasLocationID:   true,
		LocationID:      10,
		HasLocationName: true,
		LocationName:    " floor ",
		HasStatus:       true,
		Status:          " online ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.DeviceName != "new" || updated.DeviceType != "meter" || updated.RoomName != "hall" || updated.LocationID != 10 || updated.Status != "online" {
		t.Fatalf("unexpected updated device: %+v", updated)
	}
	if deviceRepo.lastUpdateDevice.UpdatedAt <= base.UpdatedAt {
		t.Fatalf("expected UpdatedAt to change: %+v", deviceRepo.lastUpdateDevice)
	}
}

func TestListDevices(t *testing.T) {
	svc := newDeviceServiceForTest(&mockUserRepo{}, &mockDeviceRepo{}, &mockThresholdRepo{})
	if _, err := svc.ListDevices(context.Background(), 0); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation, got %v", err)
	}

	for _, tc := range []struct {
		userErr   error
		listErr   error
		expected  error
		stageList bool
	}{
		{userErr: storage.ErrNotFound, expected: service.ErrNotFound},
		{userErr: storage.ErrConflict, expected: service.ErrConflict},
		{userErr: errors.New("boom"), expected: errors.New("boom")},
		{listErr: storage.ErrNotFound, expected: service.ErrNotFound, stageList: true},
		{listErr: storage.ErrConflict, expected: service.ErrConflict, stageList: true},
		{listErr: errors.New("list boom"), expected: errors.New("list boom"), stageList: true},
	} {
		userRepo := &mockUserRepo{
			getUserByIDFn: func(context.Context, int64) (storage.User, error) {
				if tc.stageList {
					return storage.User{UserID: 1}, nil
				}
				return storage.User{}, tc.userErr
			},
		}
		deviceRepo := &mockDeviceRepo{
			listDevicesByUserFn: func(context.Context, int64) ([]storage.Device, error) {
				return nil, tc.listErr
			},
		}
		svc = newDeviceServiceForTest(userRepo, deviceRepo, &mockThresholdRepo{})
		_, err := svc.ListDevices(context.Background(), 1)
		if tc.expected.Error() == "boom" || tc.expected.Error() == "list boom" {
			if err == nil || err.Error() != tc.expected.Error() {
				t.Fatalf("expected %v, got %v", tc.expected, err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	svc = newDeviceServiceForTest(
		&mockUserRepo{getUserByIDFn: func(context.Context, int64) (storage.User, error) { return storage.User{UserID: 1}, nil }},
		&mockDeviceRepo{listDevicesByUserFn: func(context.Context, int64) ([]storage.Device, error) {
			return []storage.Device{
				{DeviceID: 2, AddedAt: 100},
				{DeviceID: 1, AddedAt: 100},
				{DeviceID: 3, AddedAt: 200},
			}, nil
		}},
		&mockThresholdRepo{},
	)
	items, err := svc.ListDevices(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 || items[0].DeviceID != 3 || items[1].DeviceID != 1 || items[2].DeviceID != 2 {
		t.Fatalf("unexpected sorting: %+v", items)
	}
}

func TestGetAndPutThresholds(t *testing.T) {
	deviceRepo := &mockDeviceRepo{
		getDeviceFn: func(context.Context, int64) (storage.Device, error) {
			return storage.Device{DeviceID: 5, UserID: 1}, nil
		},
	}
	thresholdRepo := &mockThresholdRepo{}
	svc := newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, thresholdRepo)

	if _, err := svc.GetThresholds(context.Background(), 0, 5); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation from ownedDevice, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		thresholdRepo = &mockThresholdRepo{
			getThresholdsFn: func(context.Context, int64) ([]storage.DeviceThreshold, error) { return nil, tc.err },
		}
		svc = newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, thresholdRepo)
		_, err := svc.GetThresholds(context.Background(), 1, 5)
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

	thresholdRepo = &mockThresholdRepo{
		getThresholdsFn: func(context.Context, int64) ([]storage.DeviceThreshold, error) {
			return []storage.DeviceThreshold{
				{MetricType: "zeta", MaxValue: 10, HasMaxValue: true},
				{MetricType: "alpha", MinValue: 1, HasMinValue: true},
			}, nil
		},
	}
	svc = newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, thresholdRepo)
	items, err := svc.GetThresholds(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || items[0].MetricType != "alpha" || items[1].MetricType != "zeta" {
		t.Fatalf("unexpected threshold ordering: %+v", items)
	}

	if _, err := svc.PutThresholds(context.Background(), domain.ThresholdsUpsert{UserID: 1, DeviceID: 5}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation for empty thresholds, got %v", err)
	}

	for _, cmd := range []domain.ThresholdsUpsert{
		{
			UserID: 1, DeviceID: 5,
			Thresholds: []domain.ThresholdPatch{{MetricType: " ", HasMinValue: true, MinValue: 1}},
		},
		{
			UserID: 1, DeviceID: 5,
			Thresholds: []domain.ThresholdPatch{{MetricType: "temp"}},
		},
		{
			UserID: 1, DeviceID: 5,
			Thresholds: []domain.ThresholdPatch{{MetricType: "temp", HasMinValue: true, MinValue: 10, HasMaxValue: true, MaxValue: 1}},
		},
	} {
		if _, err := svc.PutThresholds(context.Background(), cmd); !errors.Is(err, service.ErrValidation) {
			t.Fatalf("expected validation, got %v for cmd %+v", err, cmd)
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
		thresholdRepo = &mockThresholdRepo{
			putThresholdsFn: func(context.Context, int64, []storage.DeviceThreshold) error { return tc.err },
		}
		svc = newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, thresholdRepo)
		_, err := svc.PutThresholds(context.Background(), domain.ThresholdsUpsert{
			UserID:   1,
			DeviceID: 5,
			Thresholds: []domain.ThresholdPatch{
				{MetricType: "Temp", HasMinValue: true, MinValue: 1},
			},
		})
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

	thresholdRepo = &mockThresholdRepo{
		putThresholdsFn: func(context.Context, int64, []storage.DeviceThreshold) error { return nil },
		getThresholdsFn: func(context.Context, int64) ([]storage.DeviceThreshold, error) {
			return []storage.DeviceThreshold{
				{MetricType: "temp", MinValue: 1, HasMinValue: true, Severity: "critical"},
			}, nil
		},
	}
	svc = newDeviceServiceForTest(&mockUserRepo{}, deviceRepo, thresholdRepo)

	if _, err := svc.PutThresholds(context.Background(), domain.ThresholdsUpsert{
		UserID:   0,
		DeviceID: 5,
		Thresholds: []domain.ThresholdPatch{
			{MetricType: "temp", HasMinValue: true, MinValue: 1},
		},
	}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation from ownedDevice, got %v", err)
	}

	items, err = svc.PutThresholds(context.Background(), domain.ThresholdsUpsert{
		UserID:   1,
		DeviceID: 5,
		Thresholds: []domain.ThresholdPatch{
			{MetricType: " Temp ", HasMinValue: true, MinValue: 1},
			{MetricType: " Hum ", HasMaxValue: true, MaxValue: 90, Severity: "warn"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].MetricType != "temp" {
		t.Fatalf("unexpected thresholds result: %+v", items)
	}
	if thresholdRepo.lastPutDeviceID != 5 || thresholdRepo.lastPutThresholds[0].Severity != "critical" || thresholdRepo.lastPutThresholds[1].Severity != "warn" {
		t.Fatalf("unexpected put payload: %+v", thresholdRepo.lastPutThresholds)
	}
}
