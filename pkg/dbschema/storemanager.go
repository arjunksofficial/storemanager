package dbschema

import "time"

const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
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

type Image struct {
	BaseModel
	RequestID int64
	ImageURL  string
	StoreTID  int64
	VisitTime *time.Time
	Perimeter float64
	Status    string
}

type Store struct {
	BaseModel
	StoreID  string
	Name     string
	AreaCode int64
}