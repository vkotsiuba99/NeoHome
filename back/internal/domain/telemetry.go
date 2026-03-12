package domain

type Telemetry struct {
	TelemetryID    int64
	DeviceID       int64
	RecordedAt     int64
	MetricType     string
	MetricValue    float64
	Unit           string
	RoomName       string
	LocationName   string
	BatteryLevel   int64
	SignalStrength int64
}

type TelemetryIngest struct {
	DeviceID       int64
	RecordedAt     int64
	HasRecordedAt  bool
	MetricType     string
	MetricValue    float64
	Unit           string
	RoomName       string
	LocationName   string
	BatteryLevel   int64
	HasBattery     bool
	SignalStrength int64
	HasSignal      bool
}

type IngestResult struct {
	Telemetry Telemetry
	Alerts    []Alert
}
