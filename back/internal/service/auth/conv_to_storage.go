package auth

import (
	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

type ConvToStore struct{}

func NewStoreConv() *ConvToStore {
	return &ConvToStore{}
}

func (c *ConvToStore) DomainUserToStorage(user domain.User) storage.User {
	return storage.User{
		UserID:       user.UserID,
		Email:        user.Email,
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		Login:        user.Login,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
