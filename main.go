package main

import (
	"fmt"
	"github.com/99ridho/metrickit-backend/db"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	"github.com/99ridho/metrickit-backend/handler"
	"github.com/labstack/echo"
)

var database *sqlx.DB

func main() {
	loadConfigurationFile()
	database = db.Initialize()

	fmt.Println("Server starting...")

	router := echo.New()
	httpHandler := &handler.Handler{}

	router.Use(middleware.Logger(), middleware.Recover())

	router.GET("/hello", httpHandler.Hello)
	router.POST("/payload", httpHandler.RetrievePayload)

	err := router.Start(":8185")

	if err != nil {
		fmt.Println("Can't start server")
	}
}

func loadConfigurationFile() {
	viper.SetConfigFile("config.json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
