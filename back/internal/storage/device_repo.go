package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
)

type DeviceRepo struct {
	dbSession cassandra.Session
	log       *slog.Logger
}

func NewDeviceRepo(conn cassandra.Connection) *DeviceRepo {
	return &DeviceRepo{
		dbSession: conn.Session(),
		log:       conn.Logger(),
	}
}

func (repo *DeviceRepo) session() cassandra.Session {
	if repo == nil {
		return nil
	}
	return repo.dbSession
}

func (repo *DeviceRepo) logger() *slog.Logger {
	if repo == nil {
		return slog.Default()
	}
	if repo.log != nil && repo.log.Handler() != nil {
		return repo.log
	}
	return slog.Default()
}

func (repo *DeviceRepo) CreateDevice(ctx context.Context, device Device) error {
	repo.logger().Info("create device in storage started", "method", "CreateDevice", "device_id", device.DeviceID, "user_id", device.UserID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO devices_by_id (device_id, user_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		device.DeviceID,
		device.UserID,
		device.DeviceName,
		device.DeviceType,
		device.RoomName,
		device.LocationID,
		device.LocationName,
		device.Status,
		device.LastSeenAt,
		device.LastMetricAt,
		device.BatteryLevel,
		device.SignalStrength,
		device.AddedAt,
		device.UpdatedAt,
	)
	batch.Query(
		`INSERT INTO devices_by_user (user_id, device_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		device.UserID,
		device.DeviceID,
		device.DeviceName,
		device.DeviceType,
		device.RoomName,
		device.LocationID,
		device.LocationName,
		device.Status,
		device.LastSeenAt,
		device.LastMetricAt,
		device.BatteryLevel,
		device.SignalStrength,
		device.AddedAt,
		device.UpdatedAt,
	)

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("create device in storage completed", "method", "CreateDevice", "device_id", device.DeviceID, "user_id", device.UserID)
	return nil
}

func (repo *DeviceRepo) ListDevicesByUser(ctx context.Context, userID int64) ([]Device, error) {
	repo.logger().Info("list devices by user from storage started", "method", "ListDevicesByUser", "user_id", userID)

	session := repo.session()
	if session == nil {
		return nil, ErrConflict
	}

	iter := session.Query(
		`SELECT device_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at FROM devices_by_user WHERE user_id = ?`,
		userID,
	).WithContext(ctx).Iter()

	devices := make([]Device, 0, 8)
	for {
		var device Device
		if !iter.Scan(
			&device.DeviceID,
			&device.DeviceName,
			&device.DeviceType,
			&device.RoomName,
			&device.LocationID,
			&device.LocationName,
			&device.Status,
			&device.LastSeenAt,
			&device.LastMetricAt,
			&device.BatteryLevel,
			&device.SignalStrength,
			&device.AddedAt,
			&device.UpdatedAt,
		) {
			break
		}
		device.UserID = userID
		devices = append(devices, device)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	repo.logger().Info("list devices by user from storage completed", "method", "ListDevicesByUser", "user_id", userID, "devices_count", len(devices))
	return devices, nil
}

func (repo *DeviceRepo) GetDevice(ctx context.Context, deviceID int64) (Device, error) {
	repo.logger().Info("get device from storage started", "method", "GetDevice", "device_id", deviceID)

	session := repo.session()
	if session == nil {
		return Device{}, ErrConflict
	}

	var device Device
	if err := session.Query(
		`SELECT user_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at FROM devices_by_id WHERE device_id = ?`,
		deviceID,
	).WithContext(ctx).Scan(
		&device.UserID,
		&device.DeviceName,
		&device.DeviceType,
		&device.RoomName,
		&device.LocationID,
		&device.LocationName,
		&device.Status,
		&device.LastSeenAt,
		&device.LastMetricAt,
		&device.BatteryLevel,
		&device.SignalStrength,
		&device.AddedAt,
		&device.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return Device{}, ErrNotFound
		}
		return Device{}, err
	}
	device.DeviceID = deviceID

	repo.logger().Info("get device from storage completed", "method", "GetDevice", "device_id", deviceID)
	return device, nil
}

func (repo *DeviceRepo) UpdateDevice(ctx context.Context, device Device) error {
	repo.logger().Info("update device in storage started", "method", "UpdateDevice", "device_id", device.DeviceID, "user_id", device.UserID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO devices_by_id (device_id, user_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		device.DeviceID,
		device.UserID,
		device.DeviceName,
		device.DeviceType,
		device.RoomName,
		device.LocationID,
		device.LocationName,
		device.Status,
		device.LastSeenAt,
		device.LastMetricAt,
		device.BatteryLevel,
		device.SignalStrength,
		device.AddedAt,
		device.UpdatedAt,
	)
	batch.Query(
		`INSERT INTO devices_by_user (user_id, device_id, device_name, device_type, room_name, location_id, location_name, status, last_seen_at, last_metric_at, battery_level, signal_strength, added_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		device.UserID,
		device.DeviceID,
		device.DeviceName,
		device.DeviceType,
		device.RoomName,
		device.LocationID,
		device.LocationName,
		device.Status,
		device.LastSeenAt,
		device.LastMetricAt,
		device.BatteryLevel,
		device.SignalStrength,
		device.AddedAt,
		device.UpdatedAt,
	)

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("update device in storage completed", "method", "UpdateDevice", "device_id", device.DeviceID, "user_id", device.UserID)
	return nil
}

func (repo *DeviceRepo) PutThresholds(ctx context.Context, deviceID int64, thresholds []DeviceThreshold) error {
	repo.logger().Info("put device thresholds in storage started", "method", "PutThresholds", "device_id", deviceID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch)
	batch = batch.WithContext(ctx)
	batch.Query(`DELETE FROM device_thresholds_by_device WHERE device_id = ?`, deviceID)
	for _, item := range thresholds {
		batch.Query(
			`INSERT INTO device_thresholds_by_device (device_id, metric_type, min_value, max_value, severity, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			deviceID, item.MetricType, item.MinValue, item.MaxValue, item.Severity, item.UpdatedAt,
		)
	}
	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("put device thresholds in storage completed", "method", "PutThresholds", "device_id", deviceID, "threshold_count", len(thresholds))
	return nil
}

func (repo *DeviceRepo) GetThresholds(ctx context.Context, deviceID int64) ([]DeviceThreshold, error) {
	repo.logger().Info("get device thresholds from storage started", "method", "GetThresholds", "device_id", deviceID)

	session := repo.session()
	if session == nil {
		return nil, ErrConflict
	}

	iter := session.Query(
		`SELECT metric_type, min_value, max_value, severity, updated_at FROM device_thresholds_by_device WHERE device_id = ?`,
		deviceID,
	).WithContext(ctx).Iter()

	thresholds := make([]DeviceThreshold, 0, 4)
	for {
		var item DeviceThreshold
		if !iter.Scan(&item.MetricType, &item.MinValue, &item.MaxValue, &item.Severity, &item.UpdatedAt) {
			break
		}
		item.DeviceID = deviceID
		thresholds = append(thresholds, item)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	repo.logger().Info("get device thresholds from storage completed", "method", "GetThresholds", "device_id", deviceID, "threshold_count", len(thresholds))
	return thresholds, nil
}
