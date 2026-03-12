package main

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kelseyhightower/envconfig"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
	"github.com/vkotsiuba99/NeoHome/back/internal/transport/http"
	"github.com/vkotsiuba99/NeoHome/back/internal/transport/mqtt"
	"github.com/vkotsiuba99/NeoHome/back/pkg/logger"
)

type Config struct {
	DB      cassandra.Config
	HTTP    http.Config
	MQTT    mqtt.Config
	Logger  logger.Config
	Service service.Config
}

func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.DB, validation.Required),
		validation.Field(&cfg.HTTP, validation.Required),
		validation.Field(&cfg.MQTT, validation.Required),
		validation.Field(&cfg.Logger, validation.Required),
		validation.Field(&cfg.Service, validation.Required),
	)
}

func Load(ctx context.Context) (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
