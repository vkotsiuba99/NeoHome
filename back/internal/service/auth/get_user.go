package auth

import (
	"context"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func (svc *Service) GetUser(ctx context.Context, userID int64) (domain.User, error) {
	svc.log.Info("get user started", "method", "GetUser", "user_id", userID)
	if userID <= 0 {
		return domain.User{}, service.ErrValidation
	}

	User, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, mapStorageErr(err)
	}

	return svc.toDomain.StorageUserToDomain(User), nil
}
