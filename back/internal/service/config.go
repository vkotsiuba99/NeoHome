package service

import (
	"context"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	JWTSecret         []byte        `envconfig:"JWT_SECRET" required:"true"`
	TokenTTL          time.Duration `envconfig:"JWT_TTL" default:"24h"`
	MinPasswordLength int           `envconfig:"SERVICE_MIN_PASSWORD_LENGTH" default:"8"`
	DefaultSeverity   string        `envconfig:"SERVICE_DEFAULT_SEVERITY" default:"critical"`
	MaxHistoryLimit   int64         `envconfig:"SERVICE_MAX_HISTORY_LIMIT" default:"1000"`
	DefaultHistory    int64         `envconfig:"SERVICE_DEFAULT_HISTORY_LIMIT" default:"100"`
}

func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.JWTSecret, validation.Required),
		validation.Field(&cfg.TokenTTL, validation.Required),
		validation.Field(&cfg.MinPasswordLength, validation.Required),
		validation.Field(&cfg.DefaultSeverity, validation.Required),
		validation.Field(&cfg.MaxHistoryLimit, validation.Required),
		validation.Field(&cfg.DefaultHistory, validation.Required),
	)
}
