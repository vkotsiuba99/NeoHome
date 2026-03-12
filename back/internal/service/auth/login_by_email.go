package auth

import (
	"context"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/pkg/jwt"
)

func (svc *Service) LoginEmail(ctx context.Context, cmd domain.Auth) (domain.Session, error) {
	svc.log.Info("login by email started", "method", "LoginEmail")

	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	password := strings.TrimSpace(cmd.Password)
	if len(email) == 0 || len(password) == 0 {
		return domain.Session{}, service.ErrValidation
	}

	User, err := svc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.Session{}, mapStorageErr(err, service.ErrUnauthorized)
	}

	user := svc.toDomain.StorageUserToDomain(User)
	if !svc.passwordMatches(password, user.PasswordHash) {
		return domain.Session{}, service.ErrUnauthorized
	}

	token, err := jwt.NewToken(user, svc.cfg.TokenTTL, svc.cfg.JWTSecret)
	if err != nil {
		return domain.Session{}, err
	}

	session := svc.toDomain.TokenToDomainSession(token, time.Now().UTC().UnixMilli()+svc.cfg.TokenTTL.Milliseconds(), user)
	svc.log.Info("login by email completed", "method", "LoginEmail", "user_id", user.UserID)
	return session, nil
}
