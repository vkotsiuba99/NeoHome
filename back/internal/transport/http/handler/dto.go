package handler

import "encoding/json"

type (
	CreateUserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Login    string `json:"login"`
		Phone    string `json:"phone"`
	}

	LoginEmailReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Login    string `json:"login"`
		Phone    string `json:"phone"`
	}
)

type (
	UserResp struct {
		UserID    int64  `json:"userId"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Login     string `json:"login"`
		Role      string `json:"role"`
		CreatedAt int64  `json:"createdAt"`
		UpdatedAt int64  `json:"updatedAt"`
	}

	UserWrapResp struct {
		User UserResp `json:"user"`
	}

	AuthResp struct {
		AccessToken string   `json:"accessToken"`
		ExpiresAt   int64    `json:"expiresAt"`
		User        UserResp `json:"user"`
	}
)

type (
	DeviceCreateReq struct {
		DeviceName   string `json:"deviceName"`
		DeviceType   string `json:"deviceType"`
		RoomName     string `json:"roomName"`
		LocationID   int64  `json:"locationId"`
		LocationName string `json:"locationName"`
		Status       string `json:"status"`
	}

	DeviceUpdateReq struct {
		DeviceName      string `json:"deviceName"`
		HasDeviceName   bool   `json:"-"`
		DeviceType      string `json:"deviceType"`
		HasDeviceType   bool   `json:"-"`
		RoomName        string `json:"roomName"`
		HasRoomName     bool   `json:"-"`
		LocationID      int64  `json:"locationId"`
		HasLocationID   bool   `json:"-"`
		LocationName    string `json:"locationName"`
		HasLocationName bool   `json:"-"`
		Status          string `json:"status"`
		HasStatus       bool   `json:"-"`
	}

	deviceUpdateRawBody struct {
		DeviceName   *string `json:"deviceName"`
		DeviceType   *string `json:"deviceType"`
		RoomName     *string `json:"roomName"`
		LocationID   *int64  `json:"locationId"`
		LocationName *string `json:"locationName"`
		Status       *string `json:"status"`
	}

	DeviceResp struct {
		DeviceID       int64  `json:"deviceId"`
		UserID         int64  `json:"userId"`
		DeviceName     string `json:"deviceName"`
		DeviceType     string `json:"deviceType"`
		RoomName       string `json:"roomName"`
		LocationID     int64  `json:"locationId"`
		LocationName   string `json:"locationName"`
		Status         string `json:"status"`
		LastSeenAt     int64  `json:"lastSeenAt"`
		LastMetricAt   int64  `json:"lastMetricAt"`
		BatteryLevel   int64  `json:"batteryLevel"`
		SignalStrength int64  `json:"signalStrength"`
		AddedAt        int64  `json:"addedAt"`
		UpdatedAt      int64  `json:"updatedAt"`
	}

	DeviceWrapResp struct {
		Device DeviceResp `json:"device"`
	}

	DevicesResp struct {
		Devices []DeviceResp `json:"devices"`
	}
)

type (
	DeviceThresholdReq struct {
		MetricType  string  `json:"metricType"`
		MinValue    float64 `json:"minValue"`
		HasMinValue bool    `json:"-"`
		MaxValue    float64 `json:"maxValue"`
		HasMaxValue bool    `json:"-"`
		Severity    string  `json:"severity"`
	}

	ThresholdsReq struct {
		Thresholds []DeviceThresholdReq `json:"thresholds"`
	}

	thresholdBody struct {
		MetricType string   `json:"metricType"`
		MinValue   *float64 `json:"minValue"`
		MaxValue   *float64 `json:"maxValue"`
		Severity   string   `json:"severity"`
	}

	thresholdsBody struct {
		Thresholds []thresholdBody `json:"thresholds"`
	}

	ThresholdResp struct {
		MetricType  string  `json:"metricType"`
		MinValue    float64 `json:"minValue"`
		HasMinValue bool    `json:"hasMinValue"`
		MaxValue    float64 `json:"maxValue"`
		HasMaxValue bool    `json:"hasMaxValue"`
		Severity    string  `json:"severity"`
		UpdatedAt   int64   `json:"updatedAt"`
	}

	ThresholdsResp struct {
		Thresholds []ThresholdResp `json:"thresholds"`
	}
)

type (
	telemetryPayloadBody struct {
		DeviceID       int64            `json:"deviceId"`
		RecordedAt     json.RawMessage  `json:"recordedAt"`
		MetricType     string  `json:"metricType"`
		MetricValue    float64 `json:"metricValue"`
		Unit           string  `json:"unit"`
		RoomName       string  `json:"roomName"`
		LocationName   string  `json:"locationName"`
		BatteryLevel   *int64  `json:"batteryLevel"`
		SignalStrength *int64  `json:"signalStrength"`
	}

	telemetryMQTTBody struct {
		Topic   string               `json:"topic"`
		Payload telemetryPayloadBody `json:"payload"`
	}

	TelemetryIngestReq struct {
		DeviceID       int64   `json:"deviceId"`
		RecordedAt     int64   `json:"recordedAt"`
		HasRecordedAt  bool    `json:"-"`
		MetricType     string  `json:"metricType"`
		MetricValue    float64 `json:"metricValue"`
		Unit           string  `json:"unit"`
		RoomName       string  `json:"roomName"`
		LocationName   string  `json:"locationName"`
		BatteryLevel   int64   `json:"batteryLevel"`
		HasBattery     bool    `json:"-"`
		SignalStrength int64   `json:"signalStrength"`
		HasSignal      bool    `json:"-"`
	}

	TelemetryMQTTReq struct {
		Topic   string             `json:"topic"`
		Payload TelemetryIngestReq `json:"payload"`
	}

	TelemetryResp struct {
		TelemetryID    int64   `json:"telemetryId"`
		DeviceID       int64   `json:"deviceId"`
		RecordedAt     int64   `json:"recordedAt"`
		MetricType     string  `json:"metricType"`
		MetricValue    float64 `json:"metricValue"`
		Unit           string  `json:"unit"`
		RoomName       string  `json:"roomName"`
		LocationName   string  `json:"locationName"`
		BatteryLevel   int64   `json:"batteryLevel"`
		SignalStrength int64   `json:"signalStrength"`
	}

	TelemetryWrapResp struct {
		Telemetry TelemetryResp `json:"telemetry"`
	}

	TelemetryListResp struct {
		Telemetry []TelemetryResp `json:"telemetry"`
	}

	TelemetryIngestResp struct {
		Telemetry TelemetryResp `json:"telemetry"`
		Alerts    []AlertResp   `json:"alerts"`
	}
)

type (
	AlertResp struct {
		AlertID    int64  `json:"alertId"`
		LocationID int64  `json:"locationId"`
		DeviceID   int64  `json:"deviceId"`
		CreatedAt  int64  `json:"createdAt"`
		Severity   string `json:"severity"`
		Message    string `json:"message"`
		IsResolved bool   `json:"isResolved"`
		ResolvedAt int64  `json:"resolvedAt"`
	}

	AlertsResp struct {
		Alerts []AlertResp `json:"alerts"`
	}

	AlertWrapResp struct {
		Alert AlertResp `json:"alert"`
	}
)
