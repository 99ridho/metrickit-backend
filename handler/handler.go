package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/99ridho/metrickit-backend/models"
	"github.com/labstack/echo"
)

type Handler struct {
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

	fmt.Println(h.extractAppLaunchMetricsTimeToFirstDrawKey(payload))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	})
}

func (h *Handler) extractAppLaunchMetricsTimeToFirstDrawKey(payloads map[string]interface{}) []*models.AppLaunchTime {
	appLaunchMetrics := payloads["applicationLaunchMetrics"].(map[string]interface{})
	histogrammedTimeToFirstDrawKey := appLaunchMetrics["histogrammedTimeToFirstDrawKey"].(map[string]interface{})
	histogramValues := histogrammedTimeToFirstDrawKey["histogramValue"].(map[string]interface{})

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
