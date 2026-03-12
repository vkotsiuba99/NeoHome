package logger

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Level     string `envconfig:"LOG_LEVEL" default:"info"`
	Format    string `envconfig:"LOG_FORMAT" default:"text"`
	AddSource bool   `envconfig:"LOG_ADD_SOURCE" default:"true"`
}

func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.Level, validation.Required),
		validation.Field(&cfg.Format, validation.Required),
	)
}
