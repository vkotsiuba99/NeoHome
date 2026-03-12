package auth

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type mockUserRepo struct {
	createUserFn   func(ctx context.Context, user storage.User) error
	getUserByIDFn  func(ctx context.Context, userID int64) (storage.User, error)
	getByEmailFn   func(ctx context.Context, email string) (storage.User, error)
	getByLoginFn   func(ctx context.Context, login string) (storage.User, error)
	updateUserFn   func(ctx context.Context, user storage.User) error
	lastCreateUser storage.User
	lastUpdateUser storage.User
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user storage.User) error {
	m.lastCreateUser = user
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, userID int64) (storage.User, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, userID)
	}
	return storage.User{}, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (storage.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return storage.User{}, nil
}

func (m *mockUserRepo) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {
	if m.getByLoginFn != nil {
		return m.getByLoginFn(ctx, login)
	}
	return storage.User{}, nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, user storage.User) error {
	m.lastUpdateUser = user
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, user)
	}
	return nil
}

func newAuthServiceForTest(repo *mockUserRepo) *Service {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(repo, service.Config{
		JWTSecret:         []byte("test-secret"),
		TokenTTL:          time.Hour,
		MinPasswordLength: 8,
		DefaultSeverity:   "critical",
		MaxHistoryLimit:   1000,
		DefaultHistory:    100,
	}, log)
}
