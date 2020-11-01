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
	VerifyStore(context.Context, []string) (map[string]int64, error)
	CreateJob(context.Context, *model.SubmitRequest, map[string]int64) (int64, error)
	Status(context.Context) error
	Visits(context.Context) error
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

func (s *MySQLStore) VerifyStore(ctx context.Context, storeIDs []string) (storeTIDMap map[string]int64, err error) {
	storeTIDMap = make(map[string]int64)
	stores := []dbschema.Store{}
	for _, storeID := range storeIDs {
		storeTIDMap[storeID] = 0
	}
	query := fmt.Sprintf("SELECT id, store_id FROM storeDB.stores s WHERE s.store_id IN (%s);", getSliceString(storeIDs))
	fmt.Println(query)
	rows, err := s.readDB.QueryContext(ctx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		store := dbschema.Store{}
		err = rows.Scan(
			&store.ID,
			&store.StoreID,
		)
		stores = append(stores, store)
	}
	rows.Close()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	for _, store := range stores {
		storeTIDMap[store.StoreID] = store.ID
	}

	for storeID, storeTID := range storeTIDMap {
		if storeTID == 0 {
			return storeTIDMap, fmt.Errorf("invalid storeID: %s", storeID)
		}
	}
	return storeTIDMap, nil
}

func (s *MySQLStore) CreateJob(ctx context.Context, req *model.SubmitRequest, storeTIDMap map[string]int64) (jobID int64, err error) {
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
		return 0, err
	}
	jobID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, visit := range req.Visits {
		for _, imageURL := range visit.ImageURLs {
			query := `INSERT INTO storeDB.images (created_at, updated_at, request_tid, image_url, store_tid, visit_time, status) VALUES (?, ?, ?, ?, ?, ?, ?);`
			_, err := tx.Exec(
				query,
				timeNow,
				timeNow,
				jobID,
				imageURL,
				storeTIDMap[visit.StoreID],
				visit.VisitTime,
				dbschema.StatusOngoing,
			)
			if err != nil {
				return 0, err
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return jobID, nil
}
func (s *MySQLStore) Status(ctx context.Context) (err error) {
	return
}
func (s *MySQLStore) Visits(ctx context.Context) (err error) {
	return
}
