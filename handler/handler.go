package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tokopedia/metrickit-backend/helpers"
	"github.com/tokopedia/metrickit-backend/models"
	"github.com/tokopedia/metrickit-backend/services"
	"github.com/labstack/echo"
)

type Handler struct {
	LaunchMetricService services.LaunchMetricService
	SignpostService     services.SignpostService
}

func (h *Handler) RetrievePayload(c echo.Context) error {
	var payload models.RootPayload

	jsonByte, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(jsonByte, &payload)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GenericResponse{
			Header: &models.GenericResponseHeader{
				Status:   "failed",
				Messages: []string{"Can't read JSON body"},
			},
			Error: err,
		})
	}

	metadata := helpers.ExtractMetadata(payload.Data)
	appLaunchMetricsFirstDrawKey := helpers.ExtractAppLaunchColdStartMetrics(payload.Data)

	ctx := c.Request().Context()
	if ctx != nil {
		ctx = context.Background()
	}

	coldStartMetricIDs, err := h.LaunchMetricService.Store(ctx, services.LaunchTypeColdStart, appLaunchMetricsFirstDrawKey, metadata)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.GenericResponse{
			Header: &models.GenericResponseHeader{
				Status:   "failed",
				Messages: []string{err.Error()},
			},
			Error: err,
		})
	}

	signpostMetrics := helpers.ExtractSignpostMetrics(payload.Data)
	signpostIDs, err := h.SignpostService.Store(ctx, signpostMetrics, metadata)

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
		Data: map[string]interface{}{
			"cold_start_metric_ids": coldStartMetricIDs,
			"signpost_metric_ids":   signpostIDs,
		},
	})
}
