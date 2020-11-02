package impl

import (
	"context"
	"database/sql"
	"fmt"
	"storemanager/pkg/dbschema"
	"storemanager/pkg/model"
	"strings"
	"time"
)

type MySQLStoreIF interface {
	VerifyStore(context.Context, string) (bool, error)
	CreateJob(context.Context, *model.SubmitRequest) (int64, map[int64][]int64, error)
	Status(context.Context, int64) (string, error)
	FailedStores(context.Context, int64) ([]string, error)
	Visits(context.Context, int64, string, model.DateFilter) ([]model.VisitResult, error)
	GetImageByID(context.Context, int64) (string, error)
	UpdateImageByID(context.Context, int64, int, string) error
	UpdateRequestByID(context.Context, int64) error
	UpdateVisitByID(context.Context, int64) error
	GetStoreOfVisitByVisitID(context.Context, int64) (string, error)
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

func (s *MySQLStore) CreateJob(ctx context.Context, req *model.SubmitRequest) (jobID int64, visits map[int64][]int64, err error) {
	visits = make(map[int64][]int64)
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
		return 0, visits, err
	}
	jobID, err = res.LastInsertId()
	if err != nil {
		return 0, visits, err
	}
	for _, visit := range req.Visits {
		query := `INSERT INTO storeDB.visits (created_at, updated_at, request_tid, store_id, visit_time, status) VALUES (?, ?, ?, ?, ?, ?);`
		res, err := tx.Exec(
			query,
			timeNow,
			timeNow,
			jobID,
			visit.StoreID,
			visit.VisitTime,
			dbschema.StatusOngoing,
		)
		if err != nil {
			return 0, visits, err
		}
		visitTID, err := res.LastInsertId()
		if err != nil {
			return 0, visits, err
		}
		var imageIDs []int64
		for _, imageURL := range visit.ImageURLs {
			query := `INSERT INTO storeDB.images (created_at, updated_at, visit_tid, image_url, status) VALUES (?, ?, ?, ?, ?);`
			res, err := tx.Exec(
				query,
				timeNow,
				timeNow,
				visitTID,
				imageURL,
				dbschema.StatusOngoing,
			)
			if err != nil {
				return 0, visits, err
			}
			imageTID, err := res.LastInsertId()
			if err != nil {
				return 0, visits, err
			}
			imageIDs = append(imageIDs, imageTID)
		}
		visits[visitTID] = imageIDs

	}
	err = tx.Commit()
	if err != nil {
		return 0, visits, err
	}
	return jobID, visits, nil
}
func (s *MySQLStore) Status(ctx context.Context, jobID int64) (status string, err error) {
	query := fmt.Sprintf("SELECT status FROM storeDB.requests r WHERE r.id = %d;", jobID)
	fmt.Println(query)
	err = s.readDB.QueryRowContext(ctx, query).Scan(&status)
	return status, err
}

func (s *MySQLStore) FailedStores(ctx context.Context, jobID int64) (failedStoreIDs []string, err error) {

	query := fmt.Sprintf("SELECT v.store_id FROM storeDB.visits v WHERE v.request_tid = %d;", jobID)
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
func (s *MySQLStore) Visits(ctx context.Context, areaCode int64, storeID string, dateFilter model.DateFilter) (results []model.VisitResult, err error) {
	conditionPrefix := " WHERE "
	conditionBuilder := strings.Builder{}
	if areaCode != 0 {
		conditionBuilder.WriteString(conditionPrefix + fmt.Sprintf(" s.area_code =%d ", areaCode))
		conditionPrefix = " AND "
	}
	if storeID != "" {
		conditionBuilder.WriteString(conditionPrefix + fmt.Sprintf(" s.store_id = '%s' ", storeID))
		conditionPrefix = " AND "
	} //20060102150405
	if dateFilter.Start != nil {
		conditionBuilder.WriteString(conditionPrefix + fmt.Sprintf(" v.visit_time >= '%s' ", dateFilter.Start.Format("2006-01-02 15:04:05")))
		conditionPrefix = " AND "
	}
	if dateFilter.End != nil {
		conditionBuilder.WriteString(conditionPrefix + fmt.Sprintf(" v.visit_time <= '%s' ", dateFilter.End.Format("2006-01-02 15:04:05")))
	}
	query := fmt.Sprintf(`SELECT  v.id,
	s.store_id,
	s.area_code,
	s.name, 
	SUM(i.perimeter) as perimeter, 
	v.visit_time
   FROM storeDB.images i 
	   INNER JOIN 
   storeDB.visits v ON v.id = i.visit_tid 
	   INNER JOIN 
   storeDB.requests r ON r.id =v.request_tid 
	   INNER JOIN 
   storeDB.stores s ON v.store_id = s.store_id 
	   %s
	 GROUP BY v.id ;`, conditionBuilder.String())
	fmt.Println(query)
	rows, err := s.readDB.QueryContext(ctx, query)
	if err != nil {
		return
	}
	res := make(map[string]model.VisitResult)
	for rows.Next() {
		visitRes := model.VisitResult{}
		var vID, perimeter int64
		var visitTime time.Time
		err = rows.Scan(
			&vID,
			&visitRes.StoreID,
			&visitRes.Area,
			&visitRes.StoreName,
			&perimeter,
			&visitTime,
		)
		if existingValue, found := res[visitRes.StoreID]; found {
			existingValue.Datas = append(existingValue.Datas, model.VisitData{
				Date:      &visitTime,
				Perimeter: perimeter,
			})
			res[visitRes.StoreID] = existingValue
		} else {
			res[visitRes.StoreID] = model.VisitResult{
				StoreID:   visitRes.StoreID,
				Area:      visitRes.Area,
				StoreName: visitRes.StoreName,
				Datas: []model.VisitData{
					{
						Date:      &visitTime,
						Perimeter: perimeter,
					},
				},
			}
		}
	}
	rows.Close()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	for _, result := range res {
		results = append(results, result)
	}
	return
}

func (s *MySQLStore) GetImageByID(ctx context.Context, imageTID int64) (imageURL string, err error) {
	query := fmt.Sprintf("SELECT i.image_url FROM storeDB.images i WHERE i.id = %d;", imageTID)
	fmt.Println(query)
	err = s.readDB.QueryRow(query).Scan(&imageURL)
	return imageURL, err
}

func (s *MySQLStore) UpdateImageByID(ctx context.Context, imageID int64, perimeter int, status string) (err error) {
	timeNow := time.Now().UTC()
	updateStmt := fmt.Sprintf("UPDATE storeDB.images SET status = '%s', perimeter = %d, updated_at = ? WHERE id = %d;", status, perimeter, imageID)
	fmt.Println(updateStmt)
	_, err = s.masterDB.Exec(updateStmt, timeNow)
	return
}

func (s *MySQLStore) UpdateRequestByID(ctx context.Context, requestID int64) (err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM storeDB.visits v WHERE v.request_tid = %d AND (v.status = '%s' OR v.status = '%s');", requestID, dbschema.StatusFailed, dbschema.StatusOngoing)
	fmt.Println(query)
	timeNow := time.Now().UTC()
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
	updateStmt := fmt.Sprintf("UPDATE storeDB.requests SET status = '%s', updated_at = ?  WHERE id = %d;", status, requestID)
	fmt.Println(updateStmt)
	_, err = s.masterDB.Exec(updateStmt, timeNow)
	return
}

func (s *MySQLStore) UpdateVisitByID(ctx context.Context, visitTID int64) (err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM storeDB.images i WHERE i.visit_tid = %d AND (i.status = '%s' OR i.status = '%s');", visitTID, dbschema.StatusFailed, dbschema.StatusOngoing)
	fmt.Println(query)
	timeNow := time.Now().UTC()
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
	updateStmt := fmt.Sprintf("UPDATE storeDB.visits SET status = '%s', updated_at = ?  WHERE id = %d;", status, visitTID)
	fmt.Println(updateStmt)
	_, err = s.masterDB.Exec(updateStmt, timeNow)
	return
}

func (s *MySQLStore) GetStoreOfVisitByVisitID(ctx context.Context, visitTID int64) (storeID string, err error) {
	query := fmt.Sprintf("SELECT v.store_id FROM storeDB.visits v WHERE v.id = %d;", visitTID)
	fmt.Println(query)
	err = s.readDB.QueryRow(query).Scan(&storeID)
	return storeID, err
}
