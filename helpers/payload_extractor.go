package helpers

import (
	"regexp"
	"strconv"

	"github.com/tokopedia/metrickit-backend/models"
)

func ExtractMetadata(payload models.PayloadBody) *models.AppMetadata {
	return &models.AppMetadata{
		ID:          0,
		Version:     payload.AppVersion,
		BuildNumber: payload.MetaData.AppBuildVersion,
		DeviceType:  payload.MetaData.DeviceType,
		OSVersion:   payload.MetaData.OsVersion,
	}
}

func ExtractAppLaunchColdStartMetrics(payload models.PayloadBody) []*models.AppLaunchTime {
	data := payload.ApplicationLaunchMetrics.HistogrammedTimeToFirstDrawKey

	histogramValues := extractHistogramValues(data.Value)
	histogramValuesLength := len(histogramValues)

	metrics := make([]*models.AppLaunchTime, 0)
	for i := 0; i < histogramValuesLength; i++ {
		histogramValue := histogramValues[i]

		metric := &models.AppLaunchTime{
			HistogramValue: histogramValue,
		}

		metrics = append(metrics, metric)
	}

	return metrics
}

func ExtractSignpostMetrics(payload models.PayloadBody) []*models.AppSignpost {
	signpostsData := payload.SignpostMetrics

	signpostsDataLength := len(signpostsData)
	signposts := make([]*models.AppSignpost, 0)

	for i := 0; i < signpostsDataLength; i++ {
		rawSignpost := signpostsData[i]

		signpostInterval := rawSignpost.SignpostIntervalData
		histogramValues := extractHistogramValues(rawSignpost.SignpostIntervalData.HistogrammedSignpostDurations.Value)

		rawAverageMemory := signpostInterval.SignpostAverageMemory
		averageMemoryStr, _ := filterString(rawAverageMemory, "[^0-9]+")
		averageMemoryFloat, _ := strconv.ParseFloat(averageMemoryStr, 64)

		rawCumulativeCPU := signpostInterval.SignpostCumulativeCPUTime
		cumulativeCPUStr, _ := filterString(rawCumulativeCPU, "[^0-9]+")
		cumulativeCPUFloat, _ := strconv.ParseFloat(cumulativeCPUStr, 64)

		rawCumulativeLogicalWrites := signpostInterval.SignpostCumulativeLogicalWrites
		cumulativeLogicalWritesStr, _ := filterString(rawCumulativeLogicalWrites, "[^0-9]+")
		cumulativeLogicalWritesFloat, _ := strconv.ParseFloat(cumulativeLogicalWritesStr, 64)

		signpost := &models.AppSignpost{
			Name:     rawSignpost.SignpostName,
			Category: rawSignpost.SignpostCategory,
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
	histogramLength := len(rawHistogramValues)

	for i := 0; i < histogramLength; i++ {
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
