package mqtt

type telemetryPayload struct {
	DeviceID       int64   `json:"deviceId"`
	RecordedAt     []byte  `json:"recordedAt"`
	MetricType     string  `json:"metricType"`
	MetricValue    float64 `json:"metricValue"`
	Unit           string  `json:"unit"`
	RoomName       string  `json:"roomName"`
	LocationName   string  `json:"locationName"`
	BatteryLevel   *int64  `json:"batteryLevel"`
	SignalStrength *int64  `json:"signalStrength"`
}
