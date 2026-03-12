package cassandra

import (
	"context"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Hosts          []string      `envconfig:"CASSANDRA_HOSTS" default:"127.0.0.1"`
	Port           int           `envconfig:"CASSANDRA_PORT" default:"9042"`
	Keyspace       string        `envconfig:"CASSANDRA_KEYSPACE" default:"neohome"`
	Datacenter     string        `envconfig:"CASSANDRA_DATACENTER" default:"datacenter1"`
	Username       string        `envconfig:"CASSANDRA_USERNAME" default:""`
	Password       string        `envconfig:"CASSANDRA_PASSWORD" default:""`
	ConnectTimeout time.Duration `envconfig:"CASSANDRA_CONNECT_TIMEOUT" default:"5s"`
	QueryTimeout   time.Duration `envconfig:"CASSANDRA_QUERY_TIMEOUT" default:"5s"`
	Consistency    string        `envconfig:"CASSANDRA_CONSISTENCY" default:"quorum"`
}

func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.Hosts, validation.Required),
		validation.Field(&cfg.Port, validation.Required),
		validation.Field(&cfg.Keyspace, validation.Required),
		validation.Field(&cfg.Datacenter, validation.Required),
		validation.Field(&cfg.ConnectTimeout, validation.Required),
		validation.Field(&cfg.QueryTimeout, validation.Required),
		validation.Field(&cfg.Consistency, validation.Required),
	)
}
