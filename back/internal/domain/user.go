package domain

type User struct {
	UserID       int64
	Email        string
	Phone        string
	PasswordHash string
	Login        string
	Role         string
	CreatedAt    int64
	UpdatedAt    int64
}
