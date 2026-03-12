package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func TestAuthConvertersAndNew(t *testing.T) {
	if NewDomainConv() == nil {
		t.Fatal("NewDomainConv returned nil")
	}
	if NewStoreConv() == nil {
		t.Fatal("NewStoreConv returned nil")
	}

	repo := &mockUserRepo{}
	svc := newAuthServiceForTest(repo)
	if svc == nil {
		t.Fatal("New returned nil")
	}

	now := time.Now().UTC().UnixMilli()
	cmd := domain.Register{Email: "a@b.c", Phone: "111", Login: "john"}
	user := svc.toDomain.RegisterToDomainUser(cmd, 10, now, "hash", defaultRole)
	if user.UserID != 10 || user.Role != defaultRole || user.CreatedAt != now || user.UpdatedAt != now {
		t.Fatalf("unexpected NewUser result: %+v", user)
	}

	dto := svc.toStorage.DomainUserToStorage(user)
	back := svc.toDomain.StorageUserToDomain(dto)
	if back.UserID != user.UserID || back.Email != user.Email || back.PasswordHash != user.PasswordHash {
		t.Fatalf("unexpected conversion result: %+v", back)
	}

	session := svc.toDomain.TokenToDomainSession("token", 123, back)
	if session.AccessToken != "token" || session.ExpiresAt != 123 || session.User.UserID != back.UserID {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestHashPasswordAndPasswordMatches(t *testing.T) {
	svc := newAuthServiceForTest(&mockUserRepo{})

	hash, err := svc.hashPassword("password123")
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}
	if len(hash) == 0 {
		t.Fatal("hashPassword returned empty hash")
	}
	if !svc.passwordMatches("password123", hash) {
		t.Fatal("passwordMatches must be true for valid pair")
	}
	if svc.passwordMatches("wrong", hash) {
		t.Fatal("passwordMatches must be false for invalid pair")
	}

	_, err = svc.hashPassword(strings.Repeat("x", 1000))
	if err == nil {
		t.Fatal("hashPassword must fail for too long password")
	}
}

func TestCreateUser(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		svc := newAuthServiceForTest(&mockUserRepo{})
		cases := []domain.Register{
			{},
			{Email: "bad", Login: "a", Password: "password123"},
			{Email: "a@b.c", Login: "a", Password: "short"},
		}
		for _, tc := range cases {
			_, err := svc.CreateUser(context.Background(), tc)
			if !errors.Is(err, service.ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		}
	})

	t.Run("repo error mapping", func(t *testing.T) {
		cases := []struct {
			err      error
			expected error
		}{
			{err: storage.ErrNotFound, expected: service.ErrNotFound},
			{err: storage.ErrConflict, expected: service.ErrConflict},
			{err: errors.New("boom"), expected: errors.New("boom")},
		}
		for _, tc := range cases {
			repo := &mockUserRepo{
				createUserFn: func(context.Context, storage.User) error { return tc.err },
			}
			svc := newAuthServiceForTest(repo)
			_, err := svc.CreateUser(context.Background(), domain.Register{
				Email:    "a@b.c",
				Login:    "john",
				Password: "password123",
				Phone:    " 111 ",
			})
			if tc.err.Error() == "boom" {
				if err == nil || err.Error() != "boom" {
					t.Fatalf("expected boom error, got %v", err)
				}
				continue
			}
			if !errors.Is(err, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, err)
			}
		}
	})

	t.Run("hash error", func(t *testing.T) {
		svc := newAuthServiceForTest(&mockUserRepo{})
		_, err := svc.CreateUser(context.Background(), domain.Register{
			Email:    "a@b.c",
			Login:    "john",
			Password: strings.Repeat("x", 1000),
		})
		if err == nil {
			t.Fatal("expected hash error for very long password")
		}
	})

	t.Run("success normalizes fields", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthServiceForTest(repo)

		user, err := svc.CreateUser(context.Background(), domain.Register{
			Email:    "  USER@EXAMPLE.COM ",
			Login:    "  ADMIN ",
			Password: "password123",
			Phone:    " 123 ",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Email != "user@example.com" || user.Login != "admin" || user.Phone != "123" || user.Role != defaultRole {
			t.Fatalf("unexpected user normalization: %+v", user)
		}
		if repo.lastCreateUser.Email != "user@example.com" || repo.lastCreateUser.Login != "admin" {
			t.Fatalf("unexpected user DTO: %+v", repo.lastCreateUser)
		}
	})
}

func TestGetUser(t *testing.T) {
	svc := newAuthServiceForTest(&mockUserRepo{})
	if _, err := svc.GetUser(context.Background(), 0); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}

	cases := []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	}
	for _, tc := range cases {
		repo := &mockUserRepo{
			getUserByIDFn: func(context.Context, int64) (storage.User, error) { return storage.User{}, tc.err },
		}
		svc = newAuthServiceForTest(repo)
		_, err := svc.GetUser(context.Background(), 10)
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	repo := &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) {
			return storage.User{UserID: 77, Email: "u@e.c"}, nil
		},
	}
	svc = newAuthServiceForTest(repo)
	user, err := svc.GetUser(context.Background(), 77)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UserID != 77 || user.Email != "u@e.c" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestLoginEmail(t *testing.T) {
	svc := newAuthServiceForTest(&mockUserRepo{})
	if _, err := svc.LoginEmail(context.Background(), domain.Auth{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}

	repo := &mockUserRepo{
		getByEmailFn: func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrNotFound },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.LoginEmail(context.Background(), domain.Auth{Email: "a@b.c", Password: "x"}); !errors.Is(err, service.ErrUnauthorized) {
		t.Fatalf("expected unauthorized for missing user, got %v", err)
	}

	repo = &mockUserRepo{
		getByEmailFn: func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrConflict },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.LoginEmail(context.Background(), domain.Auth{Email: "a@b.c", Password: "x"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}

	repo = &mockUserRepo{
		getByEmailFn: func(context.Context, string) (storage.User, error) {
			hash, _ := newAuthServiceForTest(&mockUserRepo{}).hashPassword("password123")
			return storage.User{UserID: 7, Email: "u@e.c", PasswordHash: hash, Role: "user"}, nil
		},
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.LoginEmail(context.Background(), domain.Auth{Email: "u@e.c", Password: "bad"}); !errors.Is(err, service.ErrUnauthorized) {
		t.Fatalf("expected unauthorized for bad password, got %v", err)
	}

	repo = &mockUserRepo{
		getByEmailFn: func(context.Context, string) (storage.User, error) {
			hash, _ := newAuthServiceForTest(&mockUserRepo{}).hashPassword("password123")
			return storage.User{UserID: 7, Email: "u@e.c", PasswordHash: hash, Role: "user"}, nil
		},
	}
	svc = newAuthServiceForTest(repo)
	svc.cfg.JWTSecret = nil
	if _, err := svc.LoginEmail(context.Background(), domain.Auth{Email: "u@e.c", Password: "password123"}); err == nil {
		t.Fatal("expected jwt generation error")
	}

	repo = &mockUserRepo{
		getByEmailFn: func(context.Context, string) (storage.User, error) {
			hash, _ := newAuthServiceForTest(&mockUserRepo{}).hashPassword("password123")
			return storage.User{UserID: 7, Email: "u@e.c", PasswordHash: hash, Role: "user"}, nil
		},
	}
	svc = newAuthServiceForTest(repo)
	session, err := svc.LoginEmail(context.Background(), domain.Auth{
		Email:    " U@E.C ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(session.AccessToken) == 0 || session.User.UserID != 7 || session.ExpiresAt <= time.Now().UTC().UnixMilli() {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestUpdateUser(t *testing.T) {
	baseDTO := storage.User{
		UserID:       10,
		Email:        "old@example.com",
		Login:        "oldlogin",
		Phone:        "111",
		PasswordHash: "oldhash",
		Role:         "user",
		CreatedAt:    1,
		UpdatedAt:    1,
	}

	svc := newAuthServiceForTest(&mockUserRepo{})
	if _, err := svc.UpdateUser(context.Background(), domain.Update{}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		repo := &mockUserRepo{
			getUserByIDFn: func(context.Context, int64) (storage.User, error) { return storage.User{}, tc.err },
		}
		svc = newAuthServiceForTest(repo)
		_, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Phone: "222"})
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	repo := &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation on no-op update, got %v", err)
	}

	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Email: "bad"}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected invalid email validation, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByEmailFn: func(context.Context, string) (storage.User, error) {
			return storage.User{UserID: 999}, nil
		},
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Email: "new@example.com"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected email conflict, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByEmailFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrConflict },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Email: "new@example.com"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected conflict mapping, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByEmailFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, errors.New("boom") },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Email: "new@example.com"}); err == nil || err.Error() != "boom" {
		t.Fatalf("expected boom error, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByLoginFn:  func(context.Context, string) (storage.User, error) { return storage.User{UserID: 999}, nil },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Login: "new"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected login conflict, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByLoginFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrConflict },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Login: "new"}); !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected login conflict mapping, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByLoginFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, errors.New("boom") },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Login: "new"}); err == nil || err.Error() != "boom" {
		t.Fatalf("expected boom error, got %v", err)
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
	}
	svc = newAuthServiceForTest(repo)
	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Password: "short"}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected short password validation, got %v", err)
	}

	if _, err := svc.UpdateUser(context.Background(), domain.Update{UserID: 10, Password: strings.Repeat("x", 1000)}); err == nil {
		t.Fatal("expected hash error")
	}

	for _, tc := range []struct {
		err      error
		expected error
	}{
		{err: storage.ErrNotFound, expected: service.ErrNotFound},
		{err: storage.ErrConflict, expected: service.ErrConflict},
		{err: errors.New("boom"), expected: errors.New("boom")},
	} {
		repo = &mockUserRepo{
			getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
			getByEmailFn: func(context.Context, string) (storage.User, error) {
				return storage.User{}, storage.ErrNotFound
			},
			getByLoginFn: func(context.Context, string) (storage.User, error) {
				return storage.User{}, storage.ErrNotFound
			},
			updateUserFn: func(context.Context, storage.User) error { return tc.err },
		}
		svc = newAuthServiceForTest(repo)
		_, err := svc.UpdateUser(context.Background(), domain.Update{
			UserID:   10,
			Email:    "new@example.com",
			Login:    "newlogin",
			Phone:    "222",
			Password: "password123",
		})
		if tc.err.Error() == "boom" {
			if err == nil || err.Error() != "boom" {
				t.Fatalf("expected boom error, got %v", err)
			}
			continue
		}
		if !errors.Is(err, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, err)
		}
	}

	repo = &mockUserRepo{
		getUserByIDFn: func(context.Context, int64) (storage.User, error) { return baseDTO, nil },
		getByEmailFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrNotFound },
		getByLoginFn:  func(context.Context, string) (storage.User, error) { return storage.User{}, storage.ErrNotFound },
	}
	svc = newAuthServiceForTest(repo)
	user, err := svc.UpdateUser(context.Background(), domain.Update{
		UserID:   10,
		Email:    " NEW@EXAMPLE.COM ",
		Login:    " NEWLOGIN ",
		Phone:    " 999 ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "new@example.com" || user.Login != "newlogin" || user.Phone != "999" || user.UpdatedAt <= baseDTO.UpdatedAt {
		t.Fatalf("unexpected updated user: %+v", user)
	}
	if repo.lastUpdateUser.Email != "new@example.com" || repo.lastUpdateUser.Login != "newlogin" {
		t.Fatalf("unexpected update DTO: %+v", repo.lastUpdateUser)
	}
}
