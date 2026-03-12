package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

type ReqConv struct{ mqttTopicTemplate string }

func NewReqConv(mqttTopicTemplate string) *ReqConv {
	return &ReqConv{mqttTopicTemplate: mqttTopicTemplate}
}

func (c *ReqConv) CreateUser(ctx context.Context, svc AuthService, req CreateUserReq) (domain.User, error) {
	return svc.CreateUser(ctx, domain.Register{
		Email:    req.Email,
		Password: req.Password,
		Login:    req.Login,
		Phone:    req.Phone,
	})
}

func (c *ReqConv) LoginEmail(ctx context.Context, svc AuthService, req LoginEmailReq) (domain.Session, error) {
	return svc.LoginEmail(ctx, domain.Auth{
		Email:    req.Email,
		Password: req.Password,
	})
}

func (c *ReqConv) UpdateUser(ctx context.Context, svc AuthService, user domain.User, req UpdateUserReq) (domain.User, error) {
	return svc.UpdateUser(ctx, domain.Update{
		UserID:   user.UserID,
		Email:    req.Email,
		Phone:    req.Phone,
		Login:    req.Login,
		Password: req.Password,
	})
}

func (c *ReqConv) CreateDevice(ctx context.Context, svc DeviceService, user domain.User, req DeviceCreateReq) (domain.Device, error) {
	return svc.CreateDevice(ctx, domain.DeviceCreate{
		UserID:       user.UserID,
		DeviceName:   req.DeviceName,
		DeviceType:   req.DeviceType,
		RoomName:     req.RoomName,
		LocationID:   req.LocationID,
		LocationName: req.LocationName,
		Status:       req.Status,
	})
}

func (c *ReqConv) UpdateDevice(ctx context.Context, svc DeviceService, user domain.User, deviceID int64, req DeviceUpdateReq) (domain.Device, error) {
	return svc.UpdateDevice(ctx, domain.DeviceUpdate{
		UserID:          user.UserID,
		DeviceID:        deviceID,
		DeviceName:      req.DeviceName,
		HasDeviceName:   req.HasDeviceName,
		DeviceType:      req.DeviceType,
		HasDeviceType:   req.HasDeviceType,
		RoomName:        req.RoomName,
		HasRoomName:     req.HasRoomName,
		LocationID:      req.LocationID,
		HasLocationID:   req.HasLocationID,
		LocationName:    req.LocationName,
		HasLocationName: req.HasLocationName,
		Status:          req.Status,
		HasStatus:       req.HasStatus,
	})
}

func (c *ReqConv) ToDeviceUpdate(raw deviceUpdateRawBody) DeviceUpdateReq {
	req := DeviceUpdateReq{}

	assignString := func(src *string, dst *string, has *bool) {
		if src == nil {
			return
		}
		*dst = *src
		*has = true
	}
	assignInt64 := func(src *int64, dst *int64, has *bool) {
		if src == nil {
			return
		}
		*dst = *src
		*has = true
	}

	assignString(raw.DeviceName, &req.DeviceName, &req.HasDeviceName)
	assignString(raw.DeviceType, &req.DeviceType, &req.HasDeviceType)
	assignString(raw.RoomName, &req.RoomName, &req.HasRoomName)
	assignInt64(raw.LocationID, &req.LocationID, &req.HasLocationID)
	assignString(raw.LocationName, &req.LocationName, &req.HasLocationName)
	assignString(raw.Status, &req.Status, &req.HasStatus)

	return req
}

func (c *ReqConv) ToThresholds(body thresholdsBody) ThresholdsReq {
	thresholds := make([]DeviceThresholdReq, len(body.Thresholds))
	for i := range body.Thresholds {
		item := body.Thresholds[i]
		thresholds[i] = DeviceThresholdReq{
			MetricType: item.MetricType,
			Severity:   item.Severity,
		}
		if item.MinValue != nil {
			thresholds[i].MinValue = *item.MinValue
			thresholds[i].HasMinValue = true
		}
		if item.MaxValue != nil {
			thresholds[i].MaxValue = *item.MaxValue
			thresholds[i].HasMaxValue = true
		}
	}

	return ThresholdsReq{Thresholds: thresholds}
}

func (c *ReqConv) GetThresholds(ctx context.Context, svc DeviceService, user domain.User, deviceID int64) ([]domain.DeviceThreshold, error) {
	return svc.GetThresholds(ctx, user.UserID, deviceID)
}

func (c *ReqConv) PutThresholds(ctx context.Context, svc DeviceService, user domain.User, deviceID int64, req ThresholdsReq) ([]domain.DeviceThreshold, error) {
	thresholds := make([]domain.ThresholdPatch, len(req.Thresholds))
	for i := range req.Thresholds {
		threshold := req.Thresholds[i]
		thresholds[i] = domain.ThresholdPatch{
			MetricType:  threshold.MetricType,
			MinValue:    threshold.MinValue,
			HasMinValue: threshold.HasMinValue,
			MaxValue:    threshold.MaxValue,
			HasMaxValue: threshold.HasMaxValue,
			Severity:    threshold.Severity,
		}
	}

	return svc.PutThresholds(ctx, domain.ThresholdsUpsert{
		UserID:     user.UserID,
		DeviceID:   deviceID,
		Thresholds: thresholds,
	})
}

func (c *ReqConv) IngestTelemetry(ctx context.Context, svc TelemetryService, req TelemetryIngestReq) (domain.IngestResult, error) {
	return svc.IngestTelemetry(ctx, domain.TelemetryIngest{
		DeviceID:       req.DeviceID,
		RecordedAt:     req.RecordedAt,
		HasRecordedAt:  req.HasRecordedAt,
		MetricType:     req.MetricType,
		MetricValue:    req.MetricValue,
		Unit:           req.Unit,
		RoomName:       req.RoomName,
		LocationName:   req.LocationName,
		BatteryLevel:   req.BatteryLevel,
		HasBattery:     req.HasBattery,
		SignalStrength: req.SignalStrength,
		HasSignal:      req.HasSignal,
	})
}

func (c *ReqConv) ToTelemetryIngest(body telemetryPayloadBody) TelemetryIngestReq {
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

	recordedAt, hasRecordedAt := parseRecordedAt(body.RecordedAt)

	req := TelemetryIngestReq{
		DeviceID:      body.DeviceID,
		RecordedAt:    recordedAt,
		HasRecordedAt: hasRecordedAt,
		MetricType:    body.MetricType,
		MetricValue:   body.MetricValue,
		Unit:          body.Unit,
		RoomName:      body.RoomName,
		LocationName:  body.LocationName,
	}
	if body.BatteryLevel != nil {
		req.BatteryLevel = *body.BatteryLevel
		req.HasBattery = true
	}
	if body.SignalStrength != nil {
		req.SignalStrength = *body.SignalStrength
		req.HasSignal = true
	}

	return req
}

func (c *ReqConv) ToTelemetryMQTT(body telemetryMQTTBody) TelemetryMQTTReq {
	return TelemetryMQTTReq{
		Topic:   body.Topic,
		Payload: c.ToTelemetryIngest(body.Payload),
	}
}

func (c *ReqConv) IngestTelemetryMQTT(ctx context.Context, svc TelemetryService, req TelemetryMQTTReq) (domain.IngestResult, error) {
	cmd, err := c.toTelemetryCommand(req)
	if err != nil {
		return domain.IngestResult{}, err
	}

	return svc.IngestTelemetry(ctx, cmd)
}

func (c *ReqConv) GetLatest(ctx context.Context, svc TelemetryService, user domain.User, deviceID int64) ([]domain.Telemetry, error) {
	return svc.GetLatest(ctx, user.UserID, deviceID)
}

func (c *ReqConv) ListTelemetry(ctx context.Context, svc TelemetryService, user domain.User, deviceID int64, r *http.Request) ([]domain.Telemetry, error) {
	query, err := c.deviceTelemetryQuery(r, user, deviceID)
	if err != nil {
		return nil, err
	}

	return svc.ListTelemetry(ctx, query)
}

func (c *ReqConv) ListAlerts(ctx context.Context, svc AlertService, r *http.Request, user domain.User) ([]domain.Alert, error) {
	query, err := c.alertQuery(r, user)
	if err != nil {
		return nil, err
	}

	return svc.ListAlerts(ctx, query)
}

func (c *ReqConv) ResolveAlert(ctx context.Context, svc AlertService, user domain.User, alertID int64) (domain.Alert, error) {
	return svc.ResolveAlert(ctx, domain.AlertResolve{
		UserID:  user.UserID,
		AlertID: alertID,
	})
}

func (c *ReqConv) CurrentUser(r *http.Request) (domain.User, error) {
	const userHeaderKey = "X-User"

	rawUser := r.Header.Get(userHeaderKey)
	if len(rawUser) == 0 {
		return domain.User{}, ErrUnauthorized
	}

	var user domain.User
	if err := json.Unmarshal([]byte(rawUser), &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (c *ReqConv) PathInt64(r *http.Request, name string) (int64, error) {
	value := strings.TrimSpace(mux.Vars(r)[name])
	if len(value) == 0 {
		return 0, ErrValidation
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsedValue <= 0 {
		return 0, ErrValidation
	}

	return parsedValue, nil
}

func (c *ReqConv) alertQuery(r *http.Request, user domain.User) (domain.AlertQuery, error) {
	fromCreatedAt, hasFromCreatedAt, err := c.optionalInt64(r.URL.Query().Get("from"))
	if err != nil {
		return domain.AlertQuery{}, err
	}
	toCreatedAt, hasToCreatedAt, err := c.optionalInt64(r.URL.Query().Get("to"))
	if err != nil {
		return domain.AlertQuery{}, err
	}
	locationID, hasLocationID, err := c.optionalInt64(r.URL.Query().Get("locationId"))
	if err != nil {
		return domain.AlertQuery{}, err
	}
	if !hasLocationID {
		locationID = 0
	}

	return domain.AlertQuery{
		UserID:           user.UserID,
		LocationID:       locationID,
		FromCreatedAt:    fromCreatedAt,
		HasFromCreatedAt: hasFromCreatedAt,
		ToCreatedAt:      toCreatedAt,
		HasToCreatedAt:   hasToCreatedAt,
	}, nil
}

func (c *ReqConv) deviceTelemetryQuery(r *http.Request, user domain.User, deviceID int64) (domain.TelemetryQuery, error) {
	fromRecordedAt, hasFromRecordedAt, err := c.optionalInt64(r.URL.Query().Get("from"))
	if err != nil {
		return domain.TelemetryQuery{}, err
	}
	toRecordedAt, hasToRecordedAt, err := c.optionalInt64(r.URL.Query().Get("to"))
	if err != nil {
		return domain.TelemetryQuery{}, err
	}
	limit, err := c.limit(r.URL.Query().Get("limit"))
	if err != nil {
		return domain.TelemetryQuery{}, err
	}

	return domain.TelemetryQuery{
		UserID:            user.UserID,
		DeviceID:          deviceID,
		MetricType:        r.URL.Query().Get("metricType"),
		FromRecordedAt:    fromRecordedAt,
		HasFromRecordedAt: hasFromRecordedAt,
		ToRecordedAt:      toRecordedAt,
		HasToRecordedAt:   hasToRecordedAt,
		Limit:             limit,
	}, nil
}

func (c *ReqConv) toTelemetryCommand(req TelemetryMQTTReq) (domain.TelemetryIngest, error) {
	cmd := domain.TelemetryIngest{
		DeviceID:       req.Payload.DeviceID,
		RecordedAt:     req.Payload.RecordedAt,
		HasRecordedAt:  req.Payload.HasRecordedAt,
		MetricType:     req.Payload.MetricType,
		MetricValue:    req.Payload.MetricValue,
		Unit:           req.Payload.Unit,
		RoomName:       req.Payload.RoomName,
		LocationName:   req.Payload.LocationName,
		BatteryLevel:   req.Payload.BatteryLevel,
		HasBattery:     req.Payload.HasBattery,
		SignalStrength: req.Payload.SignalStrength,
		HasSignal:      req.Payload.HasSignal,
	}
	if cmd.DeviceID > 0 {
		return cmd, nil
	}

	deviceID, err := c.deviceIDFromTopic(req.Topic)
	if err != nil {
		return domain.TelemetryIngest{}, err
	}
	cmd.DeviceID = deviceID
	return cmd, nil
}

func (c *ReqConv) optionalInt64(rawValue string) (int64, bool, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return 0, false, nil
	}

	parsedValue, err := strconv.ParseInt(trimmedValue, 10, 64)
	if err != nil {
		return 0, false, ErrValidation
	}

	return parsedValue, true, nil
}

func (c *ReqConv) limit(rawValue string) (int64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return defaultTelemetryLimit, nil
	}

	parsedValue, err := strconv.ParseInt(trimmedValue, 10, 64)
	if err != nil || parsedValue <= 0 {
		return 0, ErrValidation
	}
	if parsedValue > maxTelemetryLimit {
		return maxTelemetryLimit, nil
	}

	return parsedValue, nil
}

func (c *ReqConv) deviceIDFromTopic(topic string) (int64, error) {
	topicParts := strings.Split(strings.Trim(topic, "/"), "/")
	templateParts := strings.Split(strings.Trim(c.mqttTopicTemplate, "/"), "/")
	if len(topicParts) == 0 || len(topicParts) != len(templateParts) {
		return 0, ErrValidation
	}

	deviceIDRaw := ""
	for i := 0; i < len(templateParts); i++ {
		switch templateParts[i] {
		case "+":
			if len(deviceIDRaw) > 0 {
				return 0, ErrValidation
			}
			deviceIDRaw = strings.TrimSpace(topicParts[i])
		case "#":
			return 0, ErrValidation
		default:
			if topicParts[i] != templateParts[i] {
				return 0, ErrValidation
			}
		}
	}
	if len(deviceIDRaw) == 0 {
		return 0, ErrValidation
	}

	deviceID, err := strconv.ParseInt(deviceIDRaw, 10, 64)
	if err != nil || deviceID <= 0 {
		return 0, ErrValidation
	}

	return deviceID, nil
}
