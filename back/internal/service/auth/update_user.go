package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func (svc *Service) UpdateUser(ctx context.Context, cmd domain.Update) (domain.User, error) {
	svc.log.Info("update user started", "method", "UpdateUser", "user_id", cmd.UserID)
	if cmd.UserID <= 0 {
		return domain.User{}, service.ErrValidation
	}

	storageUser, err := svc.userRepo.GetUserByID(ctx, cmd.UserID)
	if err != nil {
		return domain.User{}, mapStorageErr(err)
	}

	normalized := svc.toDomain.UpdateToDomain(cmd)
	user := svc.toDomain.StorageUserToDomain(storageUser)
	updated := false

	if len(normalized.Email) > 0 && normalized.Email != user.Email {
		if !strings.Contains(normalized.Email, "@") {
			return domain.User{}, service.ErrValidation
		}
		existingUser, err := svc.userRepo.GetUserByEmail(ctx, normalized.Email)
		if err == nil && existingUser.UserID != user.UserID {
			return domain.User{}, service.ErrConflict
		}
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return domain.User{}, mapStorageErr(err)
		}
		user.Email = normalized.Email
		updated = true
	}

	if len(normalized.Phone) > 0 && normalized.Phone != user.Phone {
		user.Phone = normalized.Phone
		updated = true
	}

	if len(normalized.Login) > 0 && normalized.Login != user.Login {
		existingUser, err := svc.userRepo.GetUserByLogin(ctx, normalized.Login)
		if err == nil && existingUser.UserID != user.UserID {
			return domain.User{}, service.ErrConflict
		}
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return domain.User{}, mapStorageErr(err)
		}
		user.Login = normalized.Login
		updated = true
	}

	if len(normalized.Password) > 0 {
		if len(normalized.Password) < svc.cfg.MinPasswordLength {
			return domain.User{}, service.ErrValidation
		}
		passwordHash, hashErr := svc.hashPassword(normalized.Password)
		if hashErr != nil {
			return domain.User{}, hashErr
		}
		user.PasswordHash = passwordHash
		updated = true
	}
	if !updated {
		return domain.User{}, service.ErrValidation
	}

	user.UpdatedAt = time.Now().UTC().UnixMilli()
	if err := svc.userRepo.UpdateUser(ctx, svc.toStorage.DomainUserToStorage(user)); err != nil {
		return domain.User{}, mapStorageErr(err)
	}

	return user, nil
}
