package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
)

func TestRepositoriesConstructorsAndLoggerFallback(t *testing.T) {
	conn := cassandra.Connection{}

	authRepo := NewAuthRepo(conn)
	if authRepo == nil {
		t.Fatal("NewAuthRepo returned nil")
	}
	if authRepo.session() != nil {
		t.Fatal("expected nil session for zero cassandra connection")
	}

	deviceRepo := NewDeviceRepo(conn)
	if deviceRepo == nil {
		t.Fatal("NewDeviceRepo returned nil")
	}
	if deviceRepo.session() != nil {
		t.Fatal("expected nil session for zero cassandra connection")
	}

	telemetryRepo := NewTelemetryRepo(conn)
	if telemetryRepo == nil {
		t.Fatal("NewTelemetryRepo returned nil")
	}
	if telemetryRepo.session() != nil {
		t.Fatal("expected nil session for zero cassandra connection")
	}

	alertRepo := NewAlertRepo(conn)
	if alertRepo == nil {
		t.Fatal("NewAlertRepo returned nil")
	}
	if alertRepo.session() != nil {
		t.Fatal("expected nil session for zero cassandra connection")
	}

	var nilAuthRepo *AuthRepo
	if nilAuthRepo.logger() == nil {
		t.Fatal("logger fallback must not be nil")
	}
}

func TestAuthRepoNilSession(t *testing.T) {
	repo := &AuthRepo{}
	ctx := context.Background()

	if err := repo.CreateUser(ctx, User{UserID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetUserByID(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetUserByEmail(ctx, "x@y.z"); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetUserByLogin(ctx, "user"); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if err := repo.UpdateUser(ctx, User{UserID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestDeviceRepoNilSession(t *testing.T) {
	repo := &DeviceRepo{}
	ctx := context.Background()

	if err := repo.CreateDevice(ctx, Device{DeviceID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.ListDevicesByUser(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetDevice(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if err := repo.UpdateDevice(ctx, Device{DeviceID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if err := repo.PutThresholds(ctx, 1, nil); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetThresholds(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestTelemetryRepoNilSession(t *testing.T) {
	repo := &TelemetryRepo{}
	ctx := context.Background()

	if err := repo.AddTelemetry(ctx, Telemetry{TelemetryID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.ListTelemetry(ctx, 1, "", 0, false, 0, false, 10); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetLatestTelemetry(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestAlertRepoNilSession(t *testing.T) {
	repo := &AlertRepo{}
	ctx := context.Background()

	if err := repo.CreateAlert(ctx, Alert{AlertID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.GetAlert(ctx, 1); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if err := repo.UpdateAlert(ctx, Alert{AlertID: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if _, err := repo.ListAlerts(ctx, 1, 0, false, 0, false); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
