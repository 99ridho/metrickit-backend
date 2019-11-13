package helpers

import (
	"regexp"
	"strconv"

	"github.com/99ridho/metrickit-backend/models"
)

func ExtractMetadata(payloads map[string]interface{}) *models.AppMetadata {
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

func ExtractAppLaunchMetrics(key string, payloads map[string]interface{}) []*models.AppLaunchTime {
	appLaunchMetrics := payloads["applicationLaunchMetrics"].(map[string]interface{})
	histogrammedTime := appLaunchMetrics[key].(map[string]interface{})
	histogramValues := histogrammedTime["histogramValue"].(map[string]interface{})

	metrics := make([]*models.AppLaunchTime, 0)
	for i := 0; i < len(histogramValues); i++ {
		idxStr := strconv.Itoa(i)
		histogramValue := histogramValues[idxStr].(map[string]interface{})

		frequency := int64(histogramValue["bucketCount"].(float64))

		rawRangeStart := histogramValue["bucketStart"].(string)
		rangeStartStr, _ := filterString(rawRangeStart, "[^0-9]+")
		rangeStartFloat, _ := strconv.ParseFloat(rangeStartStr, 64)

		rawRangeEnd := histogramValue["bucketEnd"].(string)
		rangeEndStr, _ := filterString(rawRangeEnd, "[^0-9]+")
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

func filterString(str string, regex string) (string, error) {
	reg, err := regexp.Compile(regex)

	if err != nil {
		return "", err
	}

	filteredString := reg.ReplaceAllString(str, "")

	return filteredString, nil
}
