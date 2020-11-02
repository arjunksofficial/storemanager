package model

import (
	"errors"
	"time"
)

const (
	InternalError = "Internal error. Please try after some tinme"
)

type Visit struct {
	StoreID   string     `json:"store_id" validate:"nonzero"`
	ImageURLs []string   `json:"image_url"`
	VisitTime *time.Time `json:"visit_time" validate:"nonzero"`
}

// SubmitRequest is model for submit request
// swagger:model SubmitRequest
// in: body
type SubmitRequest struct {
	Count  int64   `json:"count"`
	Visits []Visit `json:"visits"`
}

func (s SubmitRequest) Validate() error {
	if len(s.Visits) != int(s.Count) {
		return errors.New("count and number of visits are not equal")
	}
	for _, visit := range s.Visits {
		if len(visit.ImageURLs) == 0 {
			return errors.New("no image urls present")
		}
		for _, imageURL := range visit.ImageURLs {
			if imageURL == "" {
				return errors.New("empty image url present")
			}
		}
		if visit.VisitTime == nil {
			return errors.New("visit time is empty")
		}
	}
	return nil
}

// SubmitResponse is model for submit response
// swagger:model SubmitResponse
// in: body
type SubmitResponse struct {
	JobID int64  `json:"job_id,omitempty"`
	Error string `json:"error,omitempty"`
}

type StatusError struct {
	StoreID string `json:"store_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

// StatusResponse is model for submit response
// swagger:model StatusResponse
// in: body
type StatusResponse struct {
	Status string        `json:"status,omitempty"`
	JobID  int64         `json:"job_id,omitempty"`
	Errors []StatusError `json:"errors,omitempty"`
	Error  string        `json:"error,omitempty"`
}

type VisitData struct {
	Date      *time.Time `json:"date,omitempty"`
	Perimeter int64      `json:"perimeter,omitempty"`
}

type VisitResult struct {
	StoreID   string      `json:"store_id,omitempty"`
	Area      int64       `json:"area,omitempty"`
	StoreName string      `json:"store_name,omitempty"`
	Datas     []VisitData `json:"data,omitempty"`
}

// VisitResponse is model for visit response
// swagger:model VisitResponse
// in: body
type VisitResponse struct {
	Results []VisitResult `json:"results,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type DateFilter struct {
	Start *time.Time
	End   *time.Time
}
