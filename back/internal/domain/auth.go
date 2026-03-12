package domain

type Register struct {
	Email    string
	Password string
	Login    string
	Phone    string
}

type Auth struct {
	Email    string
	Password string
}

type Update struct {
	UserID   int64
	Email    string
	Phone    string
	Login    string
	Password string
}

type Session struct {
	AccessToken string
	ExpiresAt   int64
	User        User
}
