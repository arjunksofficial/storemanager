package main

import (
	"net/http"

	"storemanager/internal/common/config"
	"storemanager/internal/common/logger"
	"storemanager/internal/server/routes"

	"github.com/spf13/viper"
)

func init() {
	config.Init()
	logger.GetLogger()
}

func main() {
	r := routes.GetRouter()
	http.ListenAndServe(":"+viper.GetString("app_port"), r)
}
