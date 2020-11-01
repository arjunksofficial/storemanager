package service

import (
	"context"

	"storemanager/internal/common/serviceerr"
	"storemanager/internal/server/impl"
	"storemanager/pkg/model"
)

type StoreManagerIF interface {
	Submit(context.Context, *model.SubmitRequest) (int64, *serviceerr.APIError)
	Status(context.Context) *serviceerr.APIError
	Visits(context.Context) *serviceerr.APIError
}

type StoreManager struct {
	MySQLStore impl.MySQLStoreIF
}

func NewStoreManager(mysqlStore impl.MySQLStoreIF) StoreManagerIF {
	return &StoreManager{
		MySQLStore: mysqlStore,
	}
}

func (s *StoreManager) Submit(ctx context.Context, req *model.SubmitRequest) (jobID int64, sErr *serviceerr.APIError) {

	storeIDs := []string{}
	for _, visits := range req.Visits {
		storeIDs = append(storeIDs, visits.StoreID)
	}
	storeTIDMap, err := s.MySQLStore.VerifyStore(ctx, storeIDs)
	if err != nil {
		return 0, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	jobID, err = s.MySQLStore.CreateJob(ctx, req, storeTIDMap)
	if err != nil {
		return 0, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	return jobID, nil
}

func (s *StoreManager) Status(ctx context.Context) (sErr *serviceerr.APIError) {
	err := s.MySQLStore.Status(ctx)
	if err != nil {
		return &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	return nil
}

func (s *StoreManager) Visits(ctx context.Context) (sErr *serviceerr.APIError) {
	err := s.MySQLStore.Visits(ctx)
	if err != nil {
		return &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	return nil
}
