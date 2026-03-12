package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
)

type AlertRepo struct {
	dbSession cassandra.Session
	log       *slog.Logger
}

func NewAlertRepo(conn cassandra.Connection) *AlertRepo {
	return &AlertRepo{
		dbSession: conn.Session(),
		log:       conn.Logger(),
	}
}

func (repo *AlertRepo) session() cassandra.Session {
	if repo == nil {
		return nil
	}
	return repo.dbSession
}

func (repo *AlertRepo) logger() *slog.Logger {
	if repo == nil {
		return slog.Default()
	}
	if repo.log != nil && repo.log.Handler() != nil {
		return repo.log
	}
	return slog.Default()
}

func (repo *AlertRepo) CreateAlert(ctx context.Context, alert Alert) error {
	repo.logger().Info("create alert in storage started", "method", "CreateAlert", "alert_id", alert.AlertID, "device_id", alert.DeviceID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO alerts_by_id (alert_id, location_id, device_id, created_at, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.AlertID,
		alert.LocationID,
		alert.DeviceID,
		alert.CreatedAt,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)
	batch.Query(
		`INSERT INTO alerts_by_location (location_id, created_at, alert_id, device_id, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.LocationID,
		alert.CreatedAt,
		alert.AlertID,
		alert.DeviceID,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)
	batch.Query(
		`INSERT INTO alerts_by_device (device_id, created_at, alert_id, location_id, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.DeviceID,
		alert.CreatedAt,
		alert.AlertID,
		alert.LocationID,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("create alert in storage completed", "method", "CreateAlert", "alert_id", alert.AlertID)
	return nil
}

func (repo *AlertRepo) GetAlert(ctx context.Context, alertID int64) (Alert, error) {
	repo.logger().Info("get alert from storage started", "method", "GetAlert", "alert_id", alertID)

	session := repo.session()
	if session == nil {
		return Alert{}, ErrConflict
	}

	alert := Alert{AlertID: alertID}
	if err := session.Query(
		`SELECT location_id, device_id, created_at, severity, message, is_resolved, resolved_at FROM alerts_by_id WHERE alert_id = ?`,
		alertID,
	).WithContext(ctx).Scan(
		&alert.LocationID,
		&alert.DeviceID,
		&alert.CreatedAt,
		&alert.Severity,
		&alert.Message,
		&alert.IsResolved,
		&alert.ResolvedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return Alert{}, ErrNotFound
		}
		return Alert{}, err
	}

	repo.logger().Info("get alert from storage completed", "method", "GetAlert", "alert_id", alertID)
	return alert, nil
}

func (repo *AlertRepo) UpdateAlert(ctx context.Context, alert Alert) error {
	repo.logger().Info("update alert in storage started", "method", "UpdateAlert", "alert_id", alert.AlertID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	batch := session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO alerts_by_id (alert_id, location_id, device_id, created_at, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.AlertID,
		alert.LocationID,
		alert.DeviceID,
		alert.CreatedAt,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)
	batch.Query(
		`INSERT INTO alerts_by_location (location_id, created_at, alert_id, device_id, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.LocationID,
		alert.CreatedAt,
		alert.AlertID,
		alert.DeviceID,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)
	batch.Query(
		`INSERT INTO alerts_by_device (device_id, created_at, alert_id, location_id, severity, message, is_resolved, resolved_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.DeviceID,
		alert.CreatedAt,
		alert.AlertID,
		alert.LocationID,
		alert.Severity,
		alert.Message,
		alert.IsResolved,
		alert.ResolvedAt,
	)

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	repo.logger().Info("update alert in storage completed", "method", "UpdateAlert", "alert_id", alert.AlertID)
	return nil
}

func (repo *AlertRepo) ListAlerts(
	ctx context.Context,
	locationID int64,
	fromTimestamp int64,
	hasFrom bool,
	toTimestamp int64,
	hasTo bool,
) ([]Alert, error) {
	repo.logger().Info("list alerts from storage started",
		"method", "ListAlerts",
		"location_id", locationID,
		"has_from", hasFrom,
		"has_to", hasTo,
	)

	session := repo.session()
	if session == nil {
		return nil, ErrConflict
	}

	result := make([]Alert, 0, 32)

	query := `SELECT created_at, alert_id, device_id, severity, message, is_resolved, resolved_at FROM alerts_by_location WHERE location_id = ?`
	args := make([]interface{}, 0, 4)
	args = append(args, locationID)
	if hasFrom {
		query += ` AND created_at >= ?`
		args = append(args, fromTimestamp)
	}
	if hasTo {
		query += ` AND created_at <= ?`
		args = append(args, toTimestamp)
	}

	iter := session.Query(query, args...).WithContext(ctx).Iter()
	for {
		row := Alert{LocationID: locationID}
		if !iter.Scan(&row.CreatedAt, &row.AlertID, &row.DeviceID, &row.Severity, &row.Message, &row.IsResolved, &row.ResolvedAt) {
			break
		}
		result = append(result, row)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	repo.logger().Info("list alerts from storage completed", "method", "ListAlerts", "alerts_count", len(result))
	return result, nil
}
