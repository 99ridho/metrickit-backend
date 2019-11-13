package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/99ridho/metrickit-backend/helpers"
	"github.com/99ridho/metrickit-backend/models"
	"github.com/99ridho/metrickit-backend/services"
	"github.com/labstack/echo"
)

type Handler struct {
	LaunchMetricService services.LaunchMetricService
}

func (h *Handler) RetrievePayload(c echo.Context) error {
	var payloadResult map[string]interface{}

	jsonByte, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(jsonByte, &payloadResult)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GenericResponse{
			Header: &models.GenericResponseHeader{
				Status:   "failed",
				Messages: []string{"Can't read JSON body"},
			},
			Error: err,
		})
	}

	// TODO: nanti diolah
	payload, ok := payloadResult["payloads"].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusBadRequest, models.GenericResponse{
			Header: &models.GenericResponseHeader{
				Status:   "failed",
				Messages: []string{"Can't parse payload data"},
			},
			Error: errors.New("Can't parse payload data"),
		})
	}

	metadata := helpers.ExtractMetadata(payload)
	appLaunchMetricsFirstDrawKey := helpers.ExtractAppLaunchMetrics("histogrammedTimeToFirstDrawKey", payload)

	ctx := c.Request().Context()
	if ctx != nil {
		ctx = context.Background()
	}

	ids, err := h.LaunchMetricService.Store(ctx, services.LaunchTypeColdStart, appLaunchMetricsFirstDrawKey, metadata)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.GenericResponse{
			Header: &models.GenericResponseHeader{
				Status:   "failed",
				Messages: []string{err.Error()},
			},
			Error: err,
		})
	}

	return c.JSON(http.StatusOK, models.GenericResponse{
		Header: &models.GenericResponseHeader{
			Status:   "success",
			Messages: []string{},
		},
		Data: ids,
	})
}
