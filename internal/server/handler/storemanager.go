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
		responsehelper.RespondAsJSON(http.StatusInternalServerError, w, model.SubmitResponse{
			Error: sErr.Error.Error(),
		})
		return
	}
	responsehelper.RespondAsJSON(http.StatusOK, w, model.SubmitResponse{JobID: jobID})
}

func (aH *StoreManagerHandler) Status(http.ResponseWriter, *http.Request) {

}
func (aH *StoreManagerHandler) Visits(http.ResponseWriter, *http.Request) {

}
