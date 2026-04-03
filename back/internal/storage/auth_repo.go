package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
)

type AuthRepo struct {
	dbSession cassandra.Session
	log       *slog.Logger
}

func NewAuthRepo(conn cassandra.Connection) *AuthRepo {
	return &AuthRepo{
		dbSession: conn.Session(),
		log:       conn.Logger(),
	}
}

func (repo *AuthRepo) session() cassandra.Session {
	if repo == nil {
		return nil
	}
	return repo.dbSession
}

func (repo *AuthRepo) logger() *slog.Logger {
	if repo == nil {
		return slog.Default()
	}
	if repo.log != nil && repo.log.Handler() != nil {
		return repo.log
	}
	return slog.Default()
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user User) error {
	repo.logger().Info("create user in storage started", "method", "CreateUser", "user_id", user.UserID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	applied, err := session.Query(
		`INSERT INTO users_by_id (user_id, email, login, phone, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		user.UserID, user.Email, user.Login, user.Phone, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).MapScanCAS(map[string]interface{}{})
	if err != nil {
		return err
	}
	if !applied {
		return ErrConflict
	}

	applied, err = session.Query(
		`INSERT INTO users_by_email (email, user_id, login, phone, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		user.Email, user.UserID, user.Login, user.Phone, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).MapScanCAS(map[string]interface{}{})
	if err != nil {
		if cleanupErr := session.Query(`DELETE FROM users_by_id WHERE user_id = ?`, user.UserID).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_id failed", "error", cleanupErr.Error(), "user_id", user.UserID)
		}
		return err
	}
	if !applied {
		if cleanupErr := session.Query(`DELETE FROM users_by_id WHERE user_id = ?`, user.UserID).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_id failed", "error", cleanupErr.Error(), "user_id", user.UserID)
		}
		return ErrConflict
	}

	applied, err = session.Query(
		`INSERT INTO users_by_login (login, user_id, email, phone, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		user.Login, user.UserID, user.Email, user.Phone, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).MapScanCAS(map[string]interface{}{})
	if err != nil {
		if cleanupErr := session.Query(`DELETE FROM users_by_email WHERE email = ?`, user.Email).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_email failed", "error", cleanupErr.Error(), "email", user.Email)
		}
		if cleanupErr := session.Query(`DELETE FROM users_by_id WHERE user_id = ?`, user.UserID).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_id failed", "error", cleanupErr.Error(), "user_id", user.UserID)
		}
		return err
	}
	if !applied {
		if cleanupErr := session.Query(`DELETE FROM users_by_email WHERE email = ?`, user.Email).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_email failed", "error", cleanupErr.Error(), "email", user.Email)
		}
		if cleanupErr := session.Query(`DELETE FROM users_by_id WHERE user_id = ?`, user.UserID).WithContext(ctx).Exec(); cleanupErr != nil {
			repo.logger().Warn("rollback users_by_id failed", "error", cleanupErr.Error(), "user_id", user.UserID)
		}
		return ErrConflict
	}

	repo.logger().Info("create user in storage completed", "method", "CreateUser", "user_id", user.UserID)
	return nil
}

func (repo *AuthRepo) GetUserByID(ctx context.Context, userID int64) (User, error) {
	repo.logger().Info("get user by id from storage started", "method", "GetUserByID", "user_id", userID)

	session := repo.session()
	if session == nil {
		return User{}, ErrConflict
	}

	var user User
	if err := session.Query(
		`SELECT email, login, phone, role, created_at, updated_at FROM users_by_id WHERE user_id = ?`,
		userID,
	).WithContext(ctx).Scan(
		&user.Email,
		&user.Login,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	user.UserID = userID

	if err := session.Query(`SELECT password_hash FROM users_by_email WHERE email = ?`, user.Email).WithContext(ctx).Scan(&user.PasswordHash); err != nil {
		if !errors.Is(err, gocql.ErrNotFound) {
			return User{}, err
		}
	}

	repo.logger().Info("get user by id from storage completed", "method", "GetUserByID", "user_id", userID)
	return user, nil
}

func (repo *AuthRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	repo.logger().Info("get user by email from storage started", "method", "GetUserByEmail")

	session := repo.session()
	if session == nil {
		return User{}, ErrConflict
	}

	var user User
	if err := session.Query(
		`SELECT user_id, login, phone, password_hash, role, created_at, updated_at FROM users_by_email WHERE email = ?`,
		email,
	).WithContext(ctx).Scan(
		&user.UserID,
		&user.Login,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	user.Email = email

	repo.logger().Info("get user by email from storage completed", "method", "GetUserByEmail", "user_id", user.UserID)
	return user, nil
}

func (repo *AuthRepo) GetUserByLogin(ctx context.Context, login string) (User, error) {
	repo.logger().Info("get user by login from storage started", "method", "GetUserByLogin")

	session := repo.session()
	if session == nil {
		return User{}, ErrConflict
	}

	var user User
	if err := session.Query(
		`SELECT user_id, email, phone, password_hash, role, created_at, updated_at FROM users_by_login WHERE login = ?`,
		login,
	).WithContext(ctx).Scan(
		&user.UserID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	user.Login = login

	repo.logger().Info("get user by login from storage completed", "method", "GetUserByLogin", "user_id", user.UserID)
	return user, nil
}

func (repo *AuthRepo) UpdateUser(ctx context.Context, user User) error {
	repo.logger().Info("update user in storage started", "method", "UpdateUser", "user_id", user.UserID)

	session := repo.session()
	if session == nil {
		return ErrConflict
	}

	existing, err := repo.GetUserByID(ctx, user.UserID)
	if err != nil {
		return err
	}

	if err := session.Query(
		`INSERT INTO users_by_id (user_id, email, login, phone, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.UserID, user.Email, user.Login, user.Phone, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}
	if err := session.Query(
		`INSERT INTO users_by_email (email, user_id, login, phone, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email, user.UserID, user.Login, user.Phone, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}
	if err := session.Query(
		`INSERT INTO users_by_login (login, user_id, email, phone, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Login, user.UserID, user.Email, user.Phone, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}

	if existing.Email != user.Email {
		if err := session.Query(`DELETE FROM users_by_email WHERE email = ?`, existing.Email).WithContext(ctx).Exec(); err != nil {
			return err
		}
	}
	if existing.Login != user.Login {
		if err := session.Query(`DELETE FROM users_by_login WHERE login = ?`, existing.Login).WithContext(ctx).Exec(); err != nil {
			return err
		}
	}

	repo.logger().Info("update user in storage completed", "method", "UpdateUser", "user_id", user.UserID)
	return nil
}
