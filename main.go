package main

import (
	"fmt"
	"github.com/99ridho/metrickit-backend/db"
	"github.com/jmoiron/sqlx"
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
	handler := &handler.Handler{}

	router.GET("/hello", handler.Hello)
	router.POST("/payload", handler.RetrievePayload)

	router.Start(":8185")
}

func loadConfigurationFile() {
	viper.SetConfigFile("config.json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
