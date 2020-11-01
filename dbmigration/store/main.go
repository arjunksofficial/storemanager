package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"storemanager/internal/common/config"
	"storemanager/internal/common/db"

	"github.com/gocarina/gocsv"
)

type Store struct {
	AreaCode int    `csv:"AreaCode"`
	Name     string `csv:"StoreName"`
	StoreID  string `csv:"StoreID"`
}

func main() {
	config.Init()
	db, err := db.GetMySQLMasterDB()
	if err != nil {
		panic(err)
	}
	storeFile, err := os.OpenFile(filepath.Join(os.Getenv("PROJECTPATH"), "test/data/StoreMasterAssignment.csv"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer storeFile.Close()
	stores := []*Store{}
	if err := gocsv.UnmarshalFile(storeFile, &stores); err != nil {
		panic(err)
	}

	// get store count
	count := 0
	query := `SELECT COUNT(*) as count FROM storeDB.stores;`
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	row.Next()
	err = row.Scan(
		&count,
	)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println("count: ", count)
	if count != 0 {
		fmt.Println("Already migrated before")
		return
	}
	timeNow := time.Now().UTC()
	wg := sync.WaitGroup{}
	for _, store := range stores {
		go func(store *Store) {
			wg.Add(1)
			query := `INSERT INTO storeDB.stores (created_at, updated_at, store_id, name, area_code) VALUES (?, ?, ?, ?, ?);`
			res, err := db.Exec(
				query,
				timeNow,
				timeNow,
				store.StoreID,
				store.Name,
				store.AreaCode,
			)
			if err != nil {
				panic(err)
			}
			fmt.Println(res.LastInsertId())
			wg.Done()
		}(store)
	}
	wg.Wait()
	return
}
