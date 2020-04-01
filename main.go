package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/99ridho/metrickit-backend/db"
	"github.com/99ridho/metrickit-backend/services"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/99ridho/metrickit-backend/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var database *sqlx.DB
var configOnce sync.Once

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
		HangTimeService:     services.NewHangTimeService(database),
	}

	router.Use(middleware.Logger(), middleware.Recover())

	router.POST("/payload", httpHandler.RetrievePayload)

	err := router.Start(":8185")

	if err != nil {
		fmt.Println("Can't start server")
	}
}

func loadConfigurationFile() {
	configOnce.Do(func() {
		configs := []string{"/etc/metrickit-backend/config.json", "files/etc/metrickit-backend/config.json"}

		for _, config := range configs {
			viper.SetConfigFile(config)
			err := viper.ReadInConfig()

			if err != nil {
				log.Printf("fatal error config file: %s", err)
			} else {
				log.Printf("successfully load config from %s", config)
			}
		}
	})
}
