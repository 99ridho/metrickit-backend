package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99ridho/metrickit-backend/db"
	"github.com/99ridho/metrickit-backend/services"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/99ridho/metrickit-backend/handler"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

var database *sqlx.DB

func main() {
	loadConfigurationFile()
	database = db.Initialize()

	defer func() {
		err := database.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server starting...")

	router := echo.New()
	httpHandler := &handler.Handler{
		LaunchMetricService: services.NewLaunchMetricService(database),
		SignpostService:     services.NewSignpostService(database),
	}

	router.Use(middleware.Logger(), middleware.Recover())

	router.POST("/payload", httpHandler.RetrievePayload)

	err := router.Start(":8185")

	if err != nil {
		fmt.Println("Can't start server")
	}
}

func loadConfigurationFile() {
	env := os.Getenv("MB_ENV")

	configFileNames := map[string]string{
		"LOCAL":  "config.json",
		"DOCKER": "config.docker.json",
	}

	config, ok := configFileNames[env]
	if ok {
		viper.SetConfigFile(config)
	} else {
		viper.SetConfigFile("config.json")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
