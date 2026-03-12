package storage

type User struct {
	UserID       int64
	Email        string
	Phone        string
	PasswordHash string
	Login        string
	Role         string
	CreatedAt    int64
	UpdatedAt    int64
}

type Device struct {
	DeviceID       int64
	UserID         int64
	DeviceName     string
	DeviceType     string
	RoomName       string
	LocationID     int64
	LocationName   string
	Status         string
	LastSeenAt     int64
	LastMetricAt   int64
	BatteryLevel   int64
	SignalStrength int64
	AddedAt        int64
	UpdatedAt      int64
}

type DeviceThreshold struct {
	DeviceID    int64
	MetricType  string
	MinValue    float64
	HasMinValue bool
	MaxValue    float64
	HasMaxValue bool
	Severity    string
	UpdatedAt   int64
}

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

type Alert struct {
	AlertID    int64
	LocationID int64
	DeviceID   int64
	CreatedAt  int64
	Severity   string
	Message    string
	IsResolved bool
	ResolvedAt int64
}
