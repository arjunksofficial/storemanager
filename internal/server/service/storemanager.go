package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"storemanager/internal/common/serviceerr"
	"storemanager/internal/server/impl"
	"storemanager/pkg/dbschema"
	"storemanager/pkg/model"
)

type StoreManagerIF interface {
	Submit(context.Context, *model.SubmitRequest) (int64, *serviceerr.APIError)
	Status(context.Context, int64) (string, []string, *serviceerr.APIError)
	Visits(context.Context, int64, string, model.DateFilter) *serviceerr.APIError
}

type StoreManager struct {
	MySQLStore   impl.MySQLStoreIF
	ImageManager impl.ImageManagerIF
}

func NewStoreManager(mysqlStore impl.MySQLStoreIF, imageManager impl.ImageManagerIF) StoreManagerIF {
	return &StoreManager{
		MySQLStore:   mysqlStore,
		ImageManager: imageManager,
	}
}

func (s *StoreManager) Submit(ctx context.Context, req *model.SubmitRequest) (jobID int64, sErr *serviceerr.APIError) {

	storeIDs := []string{}
	for _, visits := range req.Visits {
		storeIDs = append(storeIDs, visits.StoreID)
	}
	jobID, imageTIDs, err := s.MySQLStore.CreateJob(ctx, req)
	if err != nil {
		return 0, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}

	go func(jobID int64, imageTIDs []int64) {
		ctxNew, _ := context.WithTimeout(ctx, 2*time.Minute)
		wg := sync.WaitGroup{}
		for _, imageTID := range imageTIDs {
			wg.Add(1)
			imageURL, storeID, err := s.MySQLStore.GetImageByID(ctxNew, imageTID)
			if err != nil {
				log.Println("error getting image: ", err)
				wg.Done()
				continue
			}
			validStore, err := s.MySQLStore.VerifyStore(ctxNew, storeID)
			if err != nil {
				log.Println("error verifying store: ", err)
				wg.Done()
				continue
			}
			if !validStore {
				err := s.MySQLStore.UpdateImageByID(ctxNew, imageTID, 0, dbschema.StatusFailed)
				if err != nil {
					log.Println("error updating image failed store: ", err)
				}
				wg.Done()
				continue
			}
			perimeter, err := s.ImageManager.DownloadImageAndCalculatePerimeter(ctxNew, imageURL)
			if err != nil {
				log.Println(err)
				err := s.MySQLStore.UpdateImageByID(ctxNew, imageTID, 0, dbschema.StatusFailed)
				if err != nil {
					log.Println("error updating image failed download image", err)
				}
				wg.Done()
				continue
			}
			time.Sleep(time.Duration(getRandom()) * 100 * time.Millisecond)
			err = s.MySQLStore.UpdateImageByID(ctxNew, imageTID, perimeter, dbschema.StatusCompleted)
			if err != nil {
				log.Println("error updating image completed", err)
			}
			wg.Done()
		}
		wg.Wait()
		err := s.MySQLStore.UpdateRequestByID(ctxNew, jobID)
		if err != nil {
			log.Println("error updating request completed", err)
		}
	}(jobID, imageTIDs)

	return jobID, nil
}

func (s *StoreManager) Status(ctx context.Context, jobID int64) (status string, failedStoreIDs []string, sErr *serviceerr.APIError) {
	status, err := s.MySQLStore.Status(ctx, jobID)
	fmt.Println(err)
	if err != nil && err != sql.ErrNoRows {
		return status, failedStoreIDs, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	if err == sql.ErrNoRows {
		return status, failedStoreIDs, &serviceerr.APIError{
			StatusCode: 400,
			Error:      errors.New(""),
		}
	}
	if status == dbschema.StatusCompleted || status == dbschema.StatusOngoing {
		return status, failedStoreIDs, nil
	} else if status == dbschema.StatusFailed {
		failedStoreIDs, err := s.MySQLStore.FailedStores(ctx, jobID)
		if err != nil {
			return status, failedStoreIDs, &serviceerr.APIError{
				StatusCode: 500,
				Error:      err,
			}
		}
		return status, failedStoreIDs, nil
	} else {
		return status, failedStoreIDs, &serviceerr.APIError{
			StatusCode: 500,
			Error:      errors.New("unknown status"),
		}
	}
}

func (s *StoreManager) Visits(ctx context.Context, areaCode int64, storeID string, dateFilter model.DateFilter) (sErr *serviceerr.APIError) {
	err := s.MySQLStore.Visits(ctx, areaCode, storeID, dateFilter)
	if err != nil {
		return &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	return nil
}
