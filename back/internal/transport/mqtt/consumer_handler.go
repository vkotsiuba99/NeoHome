package mqtt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (consumer *Consumer) handleMessage(_ mqtt.Client, message mqtt.Message) {
	messageCtx := consumer.ctx
	if messageCtx == nil {
		messageCtx = context.Background()
	}

	deviceID, err := consumer.CommandConv.TopicToDeviceID(message.Topic(), consumer.cfg.TopicTelemetry)
	if err != nil {
		consumer.log.Warn("mqtt message ignored: invalid topic", "topic", message.Topic(), "error", err.Error())
		return
	}

	var payload telemetryPayload
	decoder := json.NewDecoder(bytes.NewReader(message.Payload()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		consumer.log.Warn("mqtt message ignored: invalid payload", "topic", message.Topic(), "error", err.Error())
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		consumer.log.Warn("mqtt message ignored: invalid payload", "topic", message.Topic(), "error", "unexpected trailing payload")
		return
	}

	cmd := consumer.CommandConv.PayloadToTelemetryCommand(payload)
	if cmd.DeviceID <= 0 {
		cmd.DeviceID = deviceID
	}

	ingestCtx, cancel := context.WithTimeout(messageCtx, consumer.cfg.HandlerTimeout)
	defer cancel()

	if _, err := consumer.service.IngestTelemetry(ingestCtx, cmd); err != nil {
		consumer.log.Error("mqtt ingest failed", "topic", message.Topic(), "error", err.Error())
		return
	}

	consumer.log.Info("mqtt telemetry ingested", "topic", message.Topic(), "device_id", cmd.DeviceID, "metric_type", cmd.MetricType)
}

func (consumer *Consumer) subscribe(client mqtt.Client) error {
	token := client.Subscribe(consumer.cfg.TopicTelemetry, consumer.cfg.QoS, consumer.handleMessage)
	if !token.WaitTimeout(consumer.cfg.ConnectTimeout) {
		return errors.New("mqtt subscribe timeout")
	}
	if err := token.Error(); err != nil {
		return err
	}

	consumer.log.Info("mqtt subscription is active", "topic", consumer.cfg.TopicTelemetry, "qos", consumer.cfg.QoS)
	return nil
}
