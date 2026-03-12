package domain

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

type DeviceCreate struct {
	UserID       int64
	DeviceName   string
	DeviceType   string
	RoomName     string
	LocationID   int64
	LocationName string
	Status       string
}

type DeviceUpdate struct {
	UserID          int64
	DeviceID        int64
	DeviceName      string
	HasDeviceName   bool
	DeviceType      string
	HasDeviceType   bool
	RoomName        string
	HasRoomName     bool
	LocationID      int64
	HasLocationID   bool
	LocationName    string
	HasLocationName bool
	Status          string
	HasStatus       bool
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

type ThresholdPatch struct {
	MetricType  string
	MinValue    float64
	HasMinValue bool
	MaxValue    float64
	HasMaxValue bool
	Severity    string
}

type ThresholdsUpsert struct {
	UserID     int64
	DeviceID   int64
	Thresholds []ThresholdPatch
}

type TelemetryQuery struct {
	UserID            int64
	DeviceID          int64
	MetricType        string
	FromRecordedAt    int64
	HasFromRecordedAt bool
	ToRecordedAt      int64
	HasToRecordedAt   bool
	Limit             int64
}
