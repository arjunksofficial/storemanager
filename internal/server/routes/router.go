package routes

import (
	"time"

	"storemanager/internal/common/db"
	"storemanager/internal/common/logger"
	"storemanager/internal/server/handler"
	"storemanager/internal/server/impl"
	"storemanager/internal/server/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter() *chi.Mux {
	startTime := time.Now()
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	// r.Use(middleware.Timeout(60 * time.Second))
	mysqlMasterDB, err := db.GetMySQLMasterDB()
	if err != nil {
		panic(err)
	}
	mysqlReadDB, err := db.GetMySQLReadDB()
	if err != nil {
		panic(err)
	}
	logger := logger.GetLogger()
	mysqlStore := impl.NewMySQLStore(mysqlMasterDB, mysqlReadDB)
	storeManagerService := service.NewStoreManager(mysqlStore)
	sH := handler.NewStoreManager(storeManagerService, logger)
	handler.SetStoreManagerRoutes(r, sH)
	hH := handler.NewHealth(&startTime)
	handler.SetHealthRoutes(r, hH)
	return r
}
