// Package handler Store Manager API
//
// Documentation for Store Manager API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//  - text/csv
//
// swagger:meta
package handler

import (
	"net/http"
	"strconv"
	"time"

	"storemanager/internal/common/middleware"
	"storemanager/internal/common/responsehelper"
	"storemanager/internal/server/service"
	"storemanager/pkg/model"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// StoreManagerIF ..
type StoreManagerIF interface {
	Submit(http.ResponseWriter, *http.Request)
	Status(http.ResponseWriter, *http.Request)
	Visits(http.ResponseWriter, *http.Request)
}

// StoreManagerHandler is handler object for store manager APIs
type StoreManagerHandler struct {
	Service service.StoreManagerIF
	logger  *logrus.Logger
}

// NewStoreManager gives a new store manager handler
func NewStoreManager(
	storageManager service.StoreManagerIF,
	logger *logrus.Logger,
) *StoreManagerHandler {
	return &StoreManagerHandler{
		Service: storageManager,
		logger:  logger,
	}
}

// SetStoreManagerRoutes creates routes for store manager
func SetStoreManagerRoutes(r chi.Router, sH *StoreManagerHandler) {
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.With(middleware.RequestBodyParser(model.SubmitRequest{})).Post("/submit", sH.Submit)
			r.Get("/status", sH.Status)
			r.Get("/visit", sH.Visits)
		})
	})
}

// swagger:operation POST /api/submit StoreManager submit
//
// Creates jobs
//
// ---
// produces:
// - application/json
// parameters:
// - name: requestBody
//   in: body
//   description: Job request
//   required: true
//   schema:
//     "$ref": "#/definitions/SubmitRequest"
// responses:
//   '201':
//     description: success response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"
//   '400':
//     description: error response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"
//   '500':
//     description: error response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"

// Submit accepts requests
func (sH *StoreManagerHandler) Submit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	submitReq := ctx.Value(middleware.ParsedRequest).(*model.SubmitRequest)
	err := submitReq.Validate()
	if err != nil {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.SubmitResponse{
			Error: err.Error(),
		})
		return
	}
	jobID, sErr := sH.Service.Submit(ctx, submitReq)
	if sErr != nil {
		sH.logger.Errorf("error while creating job, error: %v", sErr.Error)
		if sErr.StatusCode == 500 {
			responsehelper.RespondAsJSON(http.StatusInternalServerError, w, model.SubmitResponse{
				Error: model.InternalError,
			})
			return
		}
		responsehelper.RespondAsJSON(sErr.StatusCode, w, model.SubmitResponse{
			Error: sErr.Error.Error(),
		})
		return
	}
	responsehelper.RespondAsJSON(http.StatusCreated, w, model.SubmitResponse{JobID: jobID})
}

// swagger:operation GET /api/status StoreManager status
//
// Get status of an job
//
// ---
// produces:
// - application/json
// parameters:
// - name: jobid
//   in: path
//   description: Job ID
//   required: true
// responses:
//   '200':
//     description: success response
//     schema:
//       "$ref": "#/definitions/StatusResponse"
//   '400':
//     description: error response
//     schema:
//       "$ref": "#/definitions/StatusResponse"
//   '500':
//     description: error response
//     schema:
//       "$ref": "#/definitions/StatusResponse"

// Status gives request status
func (sH *StoreManagerHandler) Status(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobIDString := r.URL.Query().Get("jobid")
	if jobIDString == "" {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.StatusResponse{
			Error: "jobid required",
		})
		return
	}
	jobID, err := strconv.Atoi(jobIDString)
	if err != nil {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.StatusResponse{
			Error: "enter valid jobid",
		})
		return
	}
	status, failedStoreIDs, sErr := sH.Service.Status(ctx, int64(jobID))
	if sErr != nil {
		sH.logger.Errorf("error while creating job, error: %v", sErr.Error)
		if sErr.StatusCode == 500 {
			responsehelper.RespondAsJSON(http.StatusInternalServerError, w, model.SubmitResponse{
				Error: model.InternalError,
			})
			return
		}
		responsehelper.RespondAsJSON(sErr.StatusCode, w, model.SubmitResponse{
			Error: sErr.Error.Error(),
		})
		return
	}
	failedErrors := []model.StatusError{}
	for _, failedStoreID := range failedStoreIDs {
		failedErrors = append(failedErrors, model.StatusError{
			StoreID: failedStoreID,
			Error:   failedStoreID,
		})
	}
	responsehelper.RespondAsJSON(http.StatusOK, w, model.StatusResponse{
		JobID:  int64(jobID),
		Status: status,
		Errors: failedErrors,
	})
}

// swagger:operation POST /api/submit StoreManager submit
//
// Creates jobs
//
// ---
// produces:
// - application/json
// parameters:
// - name: requestBody
//   in: body
//   description: Job request
//   required: true
//   schema:
//     "$ref": "#/definitions/SubmitRequest"
// responses:
//   '201':
//     description: success response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"
//   '400':
//     description: error response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"
//   '500':
//     description: error response
//     schema:
//       "$ref": "#/definitions/SubmitResponse"

// Submit accepts requests
func (sH *StoreManagerHandler) Visits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	areaCodeString := r.URL.Query().Get("area")
	if areaCodeString == "" {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.VisitResponse{
			Error: "area required",
		})
		return
	}
	areaCode, err := strconv.Atoi(areaCodeString)
	if err != nil {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.VisitResponse{
			Error: "enter valid jobid",
		})
		return
	}
	storeID := r.URL.Query().Get("storeid")
	if areaCodeString == "" {
		responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.VisitResponse{
			Error: "storeid required",
		})
		return
	}
	dateFilter := model.DateFilter{}
	startDate := r.URL.Query().Get("startdate")
	endDate := r.URL.Query().Get("enddate")
	if startDate != "" {
		start, err := time.Parse("20060102150405", startDate)
		if err != nil {
			responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.VisitResponse{
				Error: "date format should be <20060102150405>",
			})
			return
		}
		dateFilter.Start = &start
	}
	if endDate != "" {
		end, err := time.Parse("20060102150405", endDate)
		if err != nil {
			responsehelper.RespondAsJSON(http.StatusBadRequest, w, model.VisitResponse{
				Error: "date format should be <20060102150405>",
			})
			return
		}
		dateFilter.End = &end
	}
	sErr := sH.Service.Visits(ctx, int64(areaCode), storeID, dateFilter)
	if sErr != nil {
		sH.logger.Errorf("error while creating job, error: %v", sErr.Error)
		if sErr.StatusCode == 500 {
			responsehelper.RespondAsJSON(http.StatusInternalServerError, w, model.VisitResponse{
				Error: model.InternalError,
			})
			return
		}
		responsehelper.RespondAsJSON(http.StatusInternalServerError, w, model.VisitResponse{
			Error: sErr.Error.Error(),
		})
		return
	}
	responsehelper.RespondAsJSON(http.StatusCreated, w, model.VisitResponse{})
}
