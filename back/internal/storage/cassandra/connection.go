package cassandra

import (
	"context"
	"log/slog"

	"github.com/gocql/gocql"
)

var (
	parseConsistency = gocql.ParseConsistencyWrapper
	newClusterConfig = gocql.NewCluster
	createSession    = func(cluster *gocql.ClusterConfig) (Session, error) {
		return cluster.CreateSession()
	}
)

type Session interface {
	Query(stmt string, values ...interface{}) *gocql.Query
	NewBatch(typ gocql.BatchType) *gocql.Batch
	ExecuteBatch(batch *gocql.Batch) error
	Close()
}

type Database struct {
	cfg Config
	log slog.Logger

	connection Connection
}

type Connection struct {
	cfg     Config
	log     slog.Logger
	session Session
}

func New(ctx context.Context, cfg Config, log *slog.Logger) (*Database, error) {
	const operation = "storage.cassandra.New"
	log = log.With(
		"layer", "infrastructure",
		"component", "cassandra",
		"op", operation,
	)

	if err := cfg.ValidateWithContext(ctx); err != nil {
		log.Error("cassandra config validation failed", "error", err.Error())
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		log.Error("cassandra connection canceled", "error", err.Error())
		return nil, err
	}

	cluster := newClusterConfig(cfg.Hosts...)
	cluster.Port = cfg.Port
	cluster.Keyspace = cfg.Keyspace
	consistency, consistencyErr := parseConsistency(cfg.Consistency)
	if consistencyErr != nil {
		log.Error("invalid cassandra consistency", "value", cfg.Consistency, "error", consistencyErr.Error())
		return nil, consistencyErr
	}
	cluster.Consistency = consistency
	cluster.ConnectTimeout = cfg.ConnectTimeout
	cluster.Timeout = cfg.QueryTimeout
	cluster.DisableInitialHostLookup = true
	if len(cfg.Datacenter) > 0 {
		cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy(cfg.Datacenter)
	}
	if len(cfg.Username) > 0 {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	session, err := createSession(cluster)
	if err != nil {
		log.Error("cassandra session create failed", "error", err.Error())
		return nil, err
	}

	connection := Connection{
		cfg:     cfg,
		session: session,
		log: *log.With(
			"layer", "storage",
			"component", "cassandra_connection",
		),
	}

	log.Info("successfully initialized cassandra connection",
		"db.hosts", cfg.Hosts,
		"db.port", cfg.Port,
		"db.keyspace", cfg.Keyspace,
		"db.datacenter", cfg.Datacenter,
	)

	return &Database{
		cfg:        cfg,
		log:        *log,
		connection: connection,
	}, nil
}

func (database *Database) Close() error {
	if database == nil {
		return nil
	}

	if database.connection.session != nil {
		database.connection.session.Close()
	}

	database.log.Info("cassandra connection closed")
	return nil
}

func (database *Database) DB() Connection {
	if database == nil {
		return Connection{}
	}
	return database.connection
}

func (connection *Connection) Logger() *slog.Logger {
	if connection == nil {
		return nil
	}
	return &connection.log
}

func (connection *Connection) Session() Session {
	if connection == nil {
		return nil
	}

	return connection.session
}
