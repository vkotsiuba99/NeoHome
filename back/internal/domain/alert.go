package domain

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

type AlertQuery struct {
	UserID           int64
	LocationID       int64
	FromCreatedAt    int64
	HasFromCreatedAt bool
	ToCreatedAt      int64
	HasToCreatedAt   bool
}

type AlertResolve struct {
	UserID  int64
	AlertID int64
}
