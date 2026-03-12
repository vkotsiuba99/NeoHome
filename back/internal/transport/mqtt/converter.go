package mqtt

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

type CommandConv struct{}

func NewCommandConv() CommandConv {
	return CommandConv{}
}

func (converter CommandConv) PayloadToTelemetryCommand(payload telemetryPayload) domain.TelemetryIngest {
	parseRecordedAt := func(raw []byte) (int64, bool) {
		rawValue := bytes.TrimSpace(raw)
		if len(rawValue) == 0 {
			return 0, false
		}
		if rawValue[0] != '"' {
			if numeric, err := strconv.ParseInt(string(rawValue), 10, 64); err == nil {
				return numeric, true
			}
			return 0, true
		}

		text, err := strconv.Unquote(string(rawValue))
		if err != nil {
			return 0, true
		}
		text = strings.TrimSpace(text)

		if numeric, err := strconv.ParseInt(text, 10, 64); err == nil {
			return numeric, true
		}
		if parsedTime, err := time.Parse(time.RFC3339Nano, text); err == nil {
			return parsedTime.UTC().UnixMilli(), true
		}
		return 0, true
	}

	recordedAt, hasRecordedAt := parseRecordedAt(payload.RecordedAt)

	cmd := domain.TelemetryIngest{
		DeviceID:     payload.DeviceID,
		MetricType:   payload.MetricType,
		MetricValue:  payload.MetricValue,
		Unit:         payload.Unit,
		RoomName:     payload.RoomName,
		LocationName: payload.LocationName,
	}
	if hasRecordedAt {
		cmd.RecordedAt = recordedAt
		cmd.HasRecordedAt = true
	}
	if payload.BatteryLevel != nil {
		cmd.BatteryLevel = *payload.BatteryLevel
		cmd.HasBattery = true
	}
	if payload.SignalStrength != nil {
		cmd.SignalStrength = *payload.SignalStrength
		cmd.HasSignal = true
	}
	return cmd
}

func (converter CommandConv) TopicToDeviceID(topic string, template string) (int64, error) {
	topicParts := strings.Split(strings.Trim(topic, "/"), "/")
	templateParts := strings.Split(strings.Trim(template, "/"), "/")
	if len(topicParts) == 0 || len(topicParts) != len(templateParts) {
		return 0, errors.New("invalid topic parts")
	}

	deviceIDRaw := ""
	for i, templatePart := range templateParts {
		topicPart := strings.TrimSpace(topicParts[i])
		switch templatePart {
		case "+":
			if len(deviceIDRaw) > 0 {
				return 0, errors.New("multiple wildcard parts are not supported")
			}
			deviceIDRaw = topicPart
		case "#":
			return 0, errors.New("multi-level wildcard is not supported")
		default:
			if topicPart != templatePart {
				return 0, errors.New("topic does not match template")
			}
		}
	}

	if len(deviceIDRaw) == 0 {
		return 0, errors.New("missing device id wildcard in template")
	}

	deviceID, err := strconv.ParseInt(deviceIDRaw, 10, 64)
	if err != nil || deviceID <= 0 {
		return 0, errors.New("invalid device id")
	}

	return deviceID, nil
}
