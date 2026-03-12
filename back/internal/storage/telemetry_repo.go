package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gocql/gocql"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
)

const latestTelemetryScanLimit = 1000

type TelemetryRepo struct {
	dbSession cassandra.Session
	log       *slog.Logger
}

func NewTelemetryRepo(conn cassandra.Connection) *TelemetryRepo {
	return &TelemetryRepo{
		dbSession: conn.Session(),
		log:       conn.Logger(),
	}
}

func (repo *TelemetryRepo) session() cassandra.Session {
	if repo == nil {
		return nil
	}
	return repo.dbSession
}

func (repo *TelemetryRepo) logger() *slog.Logger {
	if repo == nil {
		return slog.Default()
	}
	if repo.log != nil && repo.log.Handler() != nil {
		return repo.log
	}
	return slog.Default()
}

func (repo *TelemetryRepo) AddTelemetry(ctx context.Context, telemetry Telemetry) error {
	repo.logger().Info("add telemetry in storage started", "method", "AddTelemetry", "device_id", telemetry.DeviceID, "telemetry_id", telemetry.TelemetryID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO telemetry_by_device (device_id, recorded_at, metric_type, telemetry_id, metric_value, unit, room_name, location_name, battery_level, signal_strength) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		telemetry.DeviceID,
		telemetry.RecordedAt,
		telemetry.MetricType,
		telemetry.TelemetryID,
		telemetry.MetricValue,
		telemetry.Unit,
		telemetry.RoomName,
		telemetry.LocationName,
		telemetry.BatteryLevel,
		telemetry.SignalStrength,
	)
	batch.Query(
		`INSERT INTO telemetry_by_device_metric (device_id, metric_type, recorded_at, telemetry_id, metric_value, unit, room_name, location_name, battery_level, signal_strength) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		telemetry.DeviceID,
		telemetry.MetricType,
		telemetry.RecordedAt,
		telemetry.TelemetryID,
		telemetry.MetricValue,
		telemetry.Unit,
		telemetry.RoomName,
		telemetry.LocationName,
		telemetry.BatteryLevel,
		telemetry.SignalStrength,
	)

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("add telemetry in storage completed", "method", "AddTelemetry", "device_id", telemetry.DeviceID, "telemetry_id", telemetry.TelemetryID)
	return nil
}

func (repo *TelemetryRepo) ListTelemetry(
	ctx context.Context,
	deviceID int64,
	metricType string,
	fromTimestamp int64,
	hasFrom bool,
	toTimestamp int64,
	hasTo bool,
	limit int64,
) ([]Telemetry, error) {
	repo.logger().Info("list device telemetry from storage started",
		"method", "ListTelemetry",
		"device_id", deviceID,
		"metric_type", metricType,
		"limit", limit,
	)

	session := repo.session()
	if session == nil {
		return nil, ErrConflict
	}

	queryBuilder := strings.Builder{}
	args := make([]interface{}, 0, 8)
	if len(metricType) > 0 {
		queryBuilder.WriteString(`SELECT recorded_at, telemetry_id, metric_value, unit, room_name, location_name, battery_level, signal_strength FROM telemetry_by_device_metric WHERE device_id = ? AND metric_type = ?`)
		args = append(args, deviceID, metricType)
	} else {
		queryBuilder.WriteString(`SELECT recorded_at, metric_type, telemetry_id, metric_value, unit, room_name, location_name, battery_level, signal_strength FROM telemetry_by_device WHERE device_id = ?`)
		args = append(args, deviceID)
	}
	if hasFrom {
		queryBuilder.WriteString(` AND recorded_at >= ?`)
		args = append(args, fromTimestamp)
	}
	if hasTo {
		queryBuilder.WriteString(` AND recorded_at <= ?`)
		args = append(args, toTimestamp)
	}
	queryBuilder.WriteString(` LIMIT ?`)
	args = append(args, limit)

	iter := session.Query(queryBuilder.String(), args...).WithContext(ctx).Iter()
	result := make([]Telemetry, 0, 32)
	for {
		row := Telemetry{DeviceID: deviceID}
		if len(metricType) > 0 {
			row.MetricType = metricType
			if !iter.Scan(
				&row.RecordedAt,
				&row.TelemetryID,
				&row.MetricValue,
				&row.Unit,
				&row.RoomName,
				&row.LocationName,
				&row.BatteryLevel,
				&row.SignalStrength,
			) {
				break
			}
		} else if !iter.Scan(
			&row.RecordedAt,
			&row.MetricType,
			&row.TelemetryID,
			&row.MetricValue,
			&row.Unit,
			&row.RoomName,
			&row.LocationName,
			&row.BatteryLevel,
			&row.SignalStrength,
		) {
			break
		}
		result = append(result, row)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	repo.logger().Info("list device telemetry from storage completed", "method", "ListTelemetry", "result_count", len(result))
	return result, nil
}

func (repo *TelemetryRepo) GetLatestTelemetry(ctx context.Context, deviceID int64) ([]Telemetry, error) {
	repo.logger().Info("get latest telemetry from storage started", "method", "GetLatestTelemetry", "device_id", deviceID)

	session := repo.session()
	if session == nil {
		return nil, ErrConflict
	}

	iter := session.Query(
		fmt.Sprintf(`SELECT recorded_at, metric_type, telemetry_id, metric_value, unit, room_name, location_name, battery_level, signal_strength FROM telemetry_by_device WHERE device_id = ? LIMIT %d`, latestTelemetryScanLimit),
		deviceID,
	).WithContext(ctx).Iter()

	result := make([]Telemetry, 0, 32)
	for {
		row := Telemetry{DeviceID: deviceID}
		if !iter.Scan(
			&row.RecordedAt,
			&row.MetricType,
			&row.TelemetryID,
			&row.MetricValue,
			&row.Unit,
			&row.RoomName,
			&row.LocationName,
			&row.BatteryLevel,
			&row.SignalStrength,
		) {
			break
		}
		result = append(result, row)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	repo.logger().Info("get latest telemetry from storage completed", "method", "GetLatestTelemetry", "device_id", deviceID, "result_count", len(result))
	return result, nil
}
