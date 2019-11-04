package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

type Handler struct {
}

func (h *Handler) Hello(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"greeting": "hello world",
		},
	})
}

func (h *Handler) RetrievePayload(c echo.Context) error {
	var payloadResult map[string]interface{}

	jsonByte, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(jsonByte, &payloadResult)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       "failed",
			"error_message": "Can't read JSON body",
		})
	}

	// TODO: nanti diolah
	payload, ok := payloadResult["payloads"].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       "failed",
			"error_message": "Can't read payloads data",
		})
	}

	fmt.Println(payload)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	})
}
