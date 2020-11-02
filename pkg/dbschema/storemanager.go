package dbschema

import "time"

const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type BaseModel struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Request struct {
	BaseModel
	Status string
}

type Visit struct {
	BaseModel
	RequestID int64
	StoreID   string
	VisitTime *time.Time
	Status    string
}

type Image struct {
	BaseModel
	VisitID   int64
	ImageURL  string
	Perimeter float64
	Status    string
}

type Store struct {
	BaseModel
	StoreID  string
	Name     string
	AreaCode int64
}
