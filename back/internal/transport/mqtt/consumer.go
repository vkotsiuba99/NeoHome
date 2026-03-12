package mqtt

import (
	"context"
	"errors"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

type TelemetryIngestor interface {
	IngestTelemetry(ctx context.Context, cmd domain.TelemetryIngest) (domain.IngestResult, error)
}

type Consumer struct {
	ctx         context.Context
	cfg         Config
	client      mqtt.Client
	service     TelemetryIngestor
	CommandConv CommandConv
	log         slog.Logger
}

func Start(ctx context.Context, cfg Config, service TelemetryIngestor, log *slog.Logger) (*Consumer, error) {
	if !cfg.Enabled {
		log.Info("mqtt consumer is disabled")
		return nil, nil
	}
	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("nil telemetry service")
	}

	consumer := &Consumer{
		ctx:         ctx,
		cfg:         cfg,
		service:     service,
		CommandConv: NewCommandConv(),
		log: *log.With("layer", "transport",
			"component", "mqtt_consumer"),
	}

	consumer.client = mqtt.NewClient(consumer.clientOptions())
	connectToken := consumer.client.Connect()
	if !connectToken.WaitTimeout(cfg.ConnectTimeout) {
		return nil, errors.New("mqtt connect timeout")
	}
	if err := connectToken.Error(); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		consumer.Close()
	}()

	return consumer, nil
}

func (consumer *Consumer) Close() {
	if consumer == nil || consumer.client == nil {
		return
	}

	if consumer.client.IsConnected() {
		consumer.client.Disconnect(consumer.cfg.DisconnectWaitMS)
	}
}

func (consumer *Consumer) clientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().
		AddBroker(consumer.cfg.Broker).
		SetClientID(consumer.cfg.ClientID).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(consumer.cfg.ReconnectInterval)
	if len(consumer.cfg.Username) > 0 {
		opts.SetUsername(consumer.cfg.Username)
		opts.SetPassword(consumer.cfg.Password)
	}

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		consumer.log.Info("mqtt connected")
		if err := consumer.subscribe(client); err != nil {
			consumer.log.Error("mqtt subscribe failed", "error", err.Error())
		}
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		if err != nil {
			consumer.log.Warn("mqtt connection lost", "error", err.Error())
			return
		}
		consumer.log.Warn("mqtt connection lost")
	})

	return opts
}
