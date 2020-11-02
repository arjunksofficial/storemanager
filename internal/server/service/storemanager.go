package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"storemanager/internal/common/serviceerr"
	"storemanager/internal/server/impl"
	"storemanager/pkg/dbschema"
	"storemanager/pkg/model"

	"github.com/sirupsen/logrus"
)

type StoreManagerIF interface {
	Submit(context.Context, *model.SubmitRequest) (int64, *serviceerr.APIError)
	Status(context.Context, int64) (string, []string, *serviceerr.APIError)
	Visits(context.Context, int64, string, model.DateFilter) ([]model.VisitResult, *serviceerr.APIError)
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

	jobID, visits, err := s.MySQLStore.CreateJob(ctx, req)
	if err != nil {
		return 0, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}

	go func(jobID int64, visits map[int64][]int64) {
		ctxNew, _ := context.WithTimeout(ctx, 10*time.Minute)
		for visitTID, imageTIDs := range visits {
			storeID, err := s.MySQLStore.GetStoreOfVisitByVisitID(ctx, visitTID)
			if err != nil {
				logrus.Println("error verifying store: ", err)
				continue
			}
			validStore, err := s.MySQLStore.VerifyStore(ctxNew, storeID)
			if err != nil {
				logrus.Println("error verifying store: ", err)
				continue
			}
			if validStore {
				for _, imageTID := range imageTIDs {
					imageURL, err := s.MySQLStore.GetImageByID(ctxNew, imageTID)
					if err != nil {
						logrus.Println("error getting image: ", err)
						continue
					}
					perimeter, err := s.ImageManager.DownloadImageAndCalculatePerimeter(imageURL)
					if err != nil {
						logrus.Println(err)
						err := s.MySQLStore.UpdateImageByID(ctxNew, imageTID, 0, dbschema.StatusFailed)
						if err != nil {
							logrus.Println("error updating image failed download image", err)
						}
						continue
					}
					time.Sleep(time.Duration(getRandom()) * 100 * time.Millisecond)
					err = s.MySQLStore.UpdateImageByID(ctxNew, imageTID, perimeter, dbschema.StatusCompleted)
					if err != nil {
						logrus.Println("error updating image completed", err)
					}
				}
			}
			err = s.MySQLStore.UpdateVisitByID(ctxNew, visitTID)
			if err != nil {
				logrus.Println("error updating request completed", err)
			}
		}
		err := s.MySQLStore.UpdateRequestByID(ctxNew, jobID)
		if err != nil {
			logrus.Println("error updating request completed", err)
		}
	}(jobID, visits)

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

func (s *StoreManager) Visits(ctx context.Context, areaCode int64, storeID string, dateFilter model.DateFilter) (visits []model.VisitResult, sErr *serviceerr.APIError) {
	visits, err := s.MySQLStore.Visits(ctx, areaCode, storeID, dateFilter)
	if err != nil {
		return nil, &serviceerr.APIError{
			StatusCode: 500,
			Error:      err,
		}
	}
	return visits, nil
}
