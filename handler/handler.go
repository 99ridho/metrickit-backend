package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

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

	metadata := h.extractMetadata(payload)
	appLaunchMetricsFirstDrawKey := h.extractAppLaunchMetrics("histogrammedTimeToFirstDrawKey", payload)

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

func (h *Handler) extractMetadata(payloads map[string]interface{}) *models.AppMetadata {
	metadata := payloads["metaData"].(map[string]interface{})
	version := payloads["appVersion"].(string)

	return &models.AppMetadata{
		ID:          0,
		Version:     version,
		BuildNumber: metadata["appBuildVersion"].(string),
		DeviceType:  metadata["deviceType"].(string),
		OSVersion:   metadata["osVersion"].(string),
	}
}

func (h *Handler) extractAppLaunchMetrics(key string, payloads map[string]interface{}) []*models.AppLaunchTime {
	appLaunchMetrics := payloads["applicationLaunchMetrics"].(map[string]interface{})
	histogrammedTime := appLaunchMetrics[key].(map[string]interface{})
	histogramValues := histogrammedTime["histogramValue"].(map[string]interface{})

	metrics := make([]*models.AppLaunchTime, 0)
	for i := 0; i < len(histogramValues); i++ {
		idxStr := strconv.Itoa(i)
		histogramValue := histogramValues[idxStr].(map[string]interface{})

		frequency := int64(histogramValue["bucketCount"].(float64))

		rawRangeStart := histogramValue["bucketStart"].(string)
		rangeStartStr, _ := h.filterString(rawRangeStart, "[^0-9]+")
		rangeStartFloat, _ := strconv.ParseFloat(rangeStartStr, 64)

		rawRangeEnd := histogramValue["bucketEnd"].(string)
		rangeEndStr, _ := h.filterString(rawRangeEnd, "[^0-9]+")
		rangeEndFloat, _ := strconv.ParseFloat(rangeEndStr, 64)

		metric := &models.AppLaunchTime{
			HistogramValue: models.HistogramValue{
				ID:         int64(i + 1),
				RangeStart: rangeStartFloat,
				RangeEnd:   rangeEndFloat,
				Frequency:  frequency,
			},
		}

		metrics = append(metrics, metric)
	}

	return metrics
}

func (h *Handler) filterString(str string, regex string) (string, error) {
	reg, err := regexp.Compile(regex)

	if err != nil {
		return "", err
	}

	filteredString := reg.ReplaceAllString(str, "")

	return filteredString, nil
}
