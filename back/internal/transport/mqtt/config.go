package mqtt

import (
	"context"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Enabled           bool          `envconfig:"MQTT_ENABLED" default:"true"`
	Broker            string        `envconfig:"MQTT_BROKER" default:"tcp://127.0.0.1:1883"`
	ClientID          string        `envconfig:"MQTT_CLIENT_ID" default:"neohome-api"`
	Username          string        `envconfig:"MQTT_USERNAME" default:""`
	Password          string        `envconfig:"MQTT_PASSWORD" default:""`
	TopicTelemetry    string        `envconfig:"MQTT_TOPIC_TELEMETRY" default:"neohome/devices/+/telemetry"`
	QoS               byte          `envconfig:"MQTT_QOS" default:"1"`
	ConnectTimeout    time.Duration `envconfig:"MQTT_CONNECT_TIMEOUT" default:"5s"`
	HandlerTimeout    time.Duration `envconfig:"MQTT_HANDLER_TIMEOUT" default:"5s"`
	ReconnectInterval time.Duration `envconfig:"MQTT_RECONNECT_INTERVAL" default:"2s"`
	DisconnectWaitMS  uint          `envconfig:"MQTT_DISCONNECT_WAIT_MS" default:"250"`
}

func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.Broker, validation.Required),
		validation.Field(&cfg.ClientID, validation.Required),
		validation.Field(&cfg.TopicTelemetry, validation.Required),
		validation.Field(&cfg.ConnectTimeout, validation.Required),
		validation.Field(&cfg.HandlerTimeout, validation.Required),
		validation.Field(&cfg.ReconnectInterval, validation.Required),
		validation.Field(&cfg.QoS, validation.Required),
	)
}
