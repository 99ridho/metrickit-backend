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
	rawHistogramValues := histogrammedTime["histogramValue"].(map[string]interface{})

	histogramValues := extractHistogramValues(rawHistogramValues)

	metrics := make([]*models.AppLaunchTime, 0)
	for i := 0; i < len(histogramValues); i++ {
		histogramValue := histogramValues[i]

		metric := &models.AppLaunchTime{
			HistogramValue: histogramValue,
		}

		metrics = append(metrics, metric)
	}

	return metrics
}

func ExtractSignpostMetrics(payloads map[string]interface{}) []*models.AppSignpost {
	signpostsData := payloads["signpostMetrics"].([]interface{})
	signposts := make([]*models.AppSignpost, 0)

	for i := 0; i < len(signpostsData); i++ {
		rawSignpost := signpostsData[i].(map[string]interface{})
		signpostInterval := rawSignpost["signpostIntervalData"].(map[string]interface{})
		histogramSignpostDurations := signpostInterval["histogrammedSignpostDurations"].(map[string]interface{})
		rawHistogramValues := histogramSignpostDurations["histogramValue"].(map[string]interface{})
		histogramValues := extractHistogramValues(rawHistogramValues)

		rawAverageMemory := signpostInterval["signpostAverageMemory"].(string)
		averageMemoryStr, _ := filterString(rawAverageMemory, "[^0-9]+")
		averageMemoryFloat, _ := strconv.ParseFloat(averageMemoryStr, 64)

		rawCumulativeCPU := signpostInterval["signpostCumulativeCPUTime"].(string)
		cumulativeCPUStr, _ := filterString(rawCumulativeCPU, "[^0-9]+")
		cumulativeCPUFloat, _ := strconv.ParseFloat(cumulativeCPUStr, 64)

		rawCumulativeLogicalWrites := signpostInterval["signpostCumulativeLogicalWrites"].(string)
		cumulativeLogicalWritesStr, _ := filterString(rawCumulativeLogicalWrites, "[^0-9]+")
		cumulativeLogicalWritesFloat, _ := strconv.ParseFloat(cumulativeLogicalWritesStr, 64)

		signpost := &models.AppSignpost{
			Name:     rawSignpost["signpostName"].(string),
			Category: rawSignpost["signpostCategory"].(string),
			SignpostInterval: models.AppSignpostInterval{
				AverageMemory:           averageMemoryFloat,
				CumulativeCPUTime:       cumulativeCPUFloat,
				CumulativeLogicalWrites: cumulativeLogicalWritesFloat,
				SignpostHistogram: models.AppSignpostHistogram{
					HistogramValues: histogramValues,
				},
			},
		}

		signposts = append(signposts, signpost)
	}

	return signposts
}

func extractHistogramValues(rawHistogramValues map[string]interface{}) []models.HistogramValue {
	histogramValues := make([]models.HistogramValue, 0)

	for i := 0; i < len(rawHistogramValues); i++ {
		idxStr := strconv.Itoa(i)
		histogramValue := rawHistogramValues[idxStr].(map[string]interface{})

		frequency := int64(histogramValue["bucketCount"].(float64))

		rawRangeStart := histogramValue["bucketStart"].(string)
		rangeStartStr, _ := filterString(rawRangeStart, "[^0-9]+")
		rangeStartFloat, _ := strconv.ParseFloat(rangeStartStr, 64)

		rawRangeEnd := histogramValue["bucketEnd"].(string)
		rangeEndStr, _ := filterString(rawRangeEnd, "[^0-9]+")
		rangeEndFloat, _ := strconv.ParseFloat(rangeEndStr, 64)

		histogramValues = append(histogramValues, models.HistogramValue{
			ID:         int64(i + 1),
			RangeStart: rangeStartFloat,
			RangeEnd:   rangeEndFloat,
			Frequency:  frequency,
		})
	}

	return histogramValues
}

func filterString(str string, regex string) (string, error) {
	reg, err := regexp.Compile(regex)

	if err != nil {
		return "", err
	}

	filteredString := reg.ReplaceAllString(str, "")

	return filteredString, nil
}
