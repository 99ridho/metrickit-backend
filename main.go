package main

import (
	"fmt"

	"github.com/99ridho/metrickit-backend/handler"
	"github.com/labstack/echo"
)

func main() {
	fmt.Println("Server starting...")

	router := echo.New()
	handler := &handler.Handler{}

	router.GET("/hello", handler.Hello)
	router.POST("/payload", handler.RetrievePayload)

	router.Start(":8185")
}
