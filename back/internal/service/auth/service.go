package auth

import (
	"context"
	"log/slog"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

const defaultRole = "user"

type UserRepo interface {
	CreateUser(ctx context.Context, user storage.User) error
	GetUserByID(ctx context.Context, userID int64) (storage.User, error)
	GetUserByEmail(ctx context.Context, email string) (storage.User, error)
	GetUserByLogin(ctx context.Context, login string) (storage.User, error)
	UpdateUser(ctx context.Context, user storage.User) error
}

type Service struct {
	userRepo  UserRepo
	toDomain  ConvToDamain
	toStorage ConvToStore
	cfg       service.Config
	log       slog.Logger
}

func New(userRepo UserRepo, cfg service.Config, log *slog.Logger) *Service {
	return &Service{
		userRepo:  userRepo,
		toDomain:  ConvToDamain{},
		toStorage: ConvToStore{},
		cfg:       cfg,
		log:       *log.With("domain", "auth"),
	}
}
