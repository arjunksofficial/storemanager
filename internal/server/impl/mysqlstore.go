package impl

import (
	"context"
	"database/sql"
	"fmt"
	"storemanager/pkg/dbschema"
	"storemanager/pkg/model"
	"time"
)

type MySQLStoreIF interface {
	VerifyStore(context.Context, string) (bool, error)
	CreateJob(context.Context, *model.SubmitRequest) (int64, []int64, error)
	Status(context.Context, int64) (string, error)
	FailedStores(context.Context, int64) ([]string, error)
	Visits(context.Context, int64, string, model.DateFilter) error
	GetImageByID(context.Context, int64) (string, string, error)
	UpdateImageByID(context.Context, int64, int, string) error
	UpdateRequestByID(context.Context, int64) error
}

type MySQLStore struct {
	masterDB *sql.DB
	readDB   *sql.DB
}

func NewMySQLStore(masterDB, readDB *sql.DB) MySQLStoreIF {
	return &MySQLStore{
		masterDB: masterDB,
		readDB:   readDB,
	}
}

func (s *MySQLStore) VerifyStore(ctx context.Context, storeID string) (valid bool, err error) {
	query := fmt.Sprintf("SELECT s.id FROM storeDB.stores s WHERE s.store_id = '%s';", storeID)
	fmt.Println(query)
	var id int64
	err = s.readDB.QueryRow(query).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (s *MySQLStore) CreateJob(ctx context.Context, req *model.SubmitRequest) (jobID int64, imageIDs []int64, err error) {
	timeNow := time.Now().UTC()
	tx, err := s.masterDB.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO storeDB.requests (created_at, updated_at, status) VALUES (?, ?, ?);`
	res, err := tx.Exec(
		query,
		timeNow,
		timeNow,
		dbschema.StatusOngoing,
	)
	if err != nil {
		return 0, imageIDs, err
	}
	jobID, err = res.LastInsertId()
	if err != nil {
		return 0, imageIDs, err
	}
	for _, visit := range req.Visits {
		for _, imageURL := range visit.ImageURLs {
			query := `INSERT INTO storeDB.images (created_at, updated_at, request_tid, image_url, store_id, visit_time, status) VALUES (?, ?, ?, ?, ?, ?, ?);`
			res, err := tx.Exec(
				query,
				timeNow,
				timeNow,
				jobID,
				imageURL,
				visit.StoreID,
				visit.VisitTime,
				dbschema.StatusOngoing,
			)
			if err != nil {
				return 0, imageIDs, err
			}
			imageTID, err := res.LastInsertId()
			if err != nil {
				return 0, imageIDs, err
			}
			imageIDs = append(imageIDs, imageTID)
		}

	}
	err = tx.Commit()
	if err != nil {
		return 0, imageIDs, err
	}
	return jobID, imageIDs, nil
}
func (s *MySQLStore) Status(ctx context.Context, jobID int64) (status string, err error) {
	query := fmt.Sprintf("SELECT status FROM storeDB.requests r WHERE r.id = %d;", jobID)
	fmt.Println(query)
	err = s.readDB.QueryRowContext(ctx, query).Scan(&status)
	return status, err
}

func (s *MySQLStore) FailedStores(ctx context.Context, jobID int64) (failedStoreIDs []string, err error) {

	query := fmt.Sprintf("SELECT i.store_id FROM storeDB.images i WHERE i.request_id = %d;", jobID)
	fmt.Println(query)
	rows, err := s.readDB.QueryContext(ctx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		failedStoreID := ""
		err = rows.Scan(
			&failedStoreID,
		)
		failedStoreIDs = append(failedStoreIDs, failedStoreID)
	}
	rows.Close()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}
func (s *MySQLStore) Visits(ctx context.Context, areaCode int64, storeID string, dateFilter model.DateFilter) (err error) {
	return
}

func (s *MySQLStore) GetImageByID(ctx context.Context, imageTID int64) (imageURL, storeID string, err error) {
	query := fmt.Sprintf("SELECT i.image_url, i.store_id FROM storeDB.images i WHERE i.id = %d;", imageTID)
	fmt.Println(query)
	err = s.readDB.QueryRow(query).Scan(&imageURL, &storeID)
	return imageURL, storeID, err
}

func (s *MySQLStore) UpdateImageByID(ctx context.Context, imageID int64, perimeter int, status string) (err error) {
	updateStmt := fmt.Sprintf("UPDATE storeDB.images SET status = '%s', perimeter = %d WHERE id = %d;", status, perimeter, imageID)
	fmt.Println(updateStmt)
	_, err = s.masterDB.Exec(updateStmt)
	return
}

func (s *MySQLStore) UpdateRequestByID(ctx context.Context, requestID int64) (err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM storeDB.images i WHERE i.request_tid = %d AND (i.status = '%s' OR i.status = '%s');", requestID, dbschema.StatusFailed, dbschema.StatusOngoing)
	fmt.Println(query)
	count := 0
	err = s.readDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	status := ""
	if count != 0 {
		status = dbschema.StatusFailed
	} else {
		status = dbschema.StatusCompleted
	}
	updateStmt := fmt.Sprintf("UPDATE storeDB.requests SET status = '%s' WHERE id = %d;", status, requestID)
	fmt.Println(updateStmt)
	_, err = s.masterDB.Exec(updateStmt)
	return
}
