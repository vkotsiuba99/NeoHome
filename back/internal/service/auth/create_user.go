package auth

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) CreateUser(ctx context.Context, cmd domain.Register) (domain.User, error) {
	svc.log.Info("create user started", "method", "CreateUser")

	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	login := strings.ToLower(strings.TrimSpace(cmd.Login))
	phone := strings.TrimSpace(cmd.Phone)
	password := strings.TrimSpace(cmd.Password)

	if len(email) == 0 || len(login) == 0 || len(password) == 0 {
		return domain.User{}, service.ErrValidation
	}
	if !strings.Contains(email, "@") {
		return domain.User{}, service.ErrValidation
	}
	if len(password) < svc.cfg.MinPasswordLength {
		return domain.User{}, service.ErrValidation
	}

	cmd.Email = email
	cmd.Login = login
	cmd.Phone = phone

	var userID int64
	for userID <= 0 {
		rawID, idErr := rand.Int(rand.Reader, big.NewInt(1<<53))
		if idErr != nil {
			return domain.User{}, idErr
		}
		userID = rawID.Int64()
	}
	now := time.Now().UTC().UnixMilli()
	passwordHash, err := svc.hashPassword(password)
	if err != nil {
		return domain.User{}, err
	}
	user := svc.toDomain.RegisterToDomainUser(cmd, userID, now, passwordHash, defaultRole)
	if err := svc.userRepo.CreateUser(ctx, svc.toStorage.DomainUserToStorage(user)); err != nil {
		return domain.User{}, mapStorageErr(err)
	}

	svc.log.Info("create user completed", "method", "CreateUser", "user_id", user.UserID)
	return user, nil
}
