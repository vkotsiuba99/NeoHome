package auth

import (
	"strings"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToDamain struct{}

func NewDomainConv() *ConvToDamain {
	return &ConvToDamain{}
}

func (c *ConvToDamain) RegisterToDomainUser(cmd domain.Register, userID int64, now int64, passwordHash string, defaultRole string) domain.User {
	return domain.User{
		UserID:       userID,
		Email:        cmd.Email,
		Phone:        cmd.Phone,
		PasswordHash: passwordHash,
		Login:        cmd.Login,
		Role:         defaultRole,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (c *ConvToDamain) StorageUserToDomain(dto storage.User) domain.User {
	return domain.User{
		UserID:       dto.UserID,
		Email:        dto.Email,
		Phone:        dto.Phone,
		PasswordHash: dto.PasswordHash,
		Login:        dto.Login,
		Role:         dto.Role,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func (c *ConvToDamain) TokenToDomainSession(token string, expiresAt int64, user domain.User) domain.Session {
	return domain.Session{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        user,
	}
}

func (c *ConvToDamain) UpdateToDomain(cmd domain.Update) domain.Update {
	return domain.Update{
		UserID:   cmd.UserID,
		Email:    strings.ToLower(strings.TrimSpace(cmd.Email)),
		Phone:    strings.TrimSpace(cmd.Phone),
		Login:    strings.ToLower(strings.TrimSpace(cmd.Login)),
		Password: strings.TrimSpace(cmd.Password),
	}
}
