package models

type HistogramValue struct {
	ID         int64
	RangeStart float64
	RangeEnd   float64
	Frequency  int64
}

type AppLaunchTime struct {
	HistogramValue
	MetadataID int64
}

type AppSignpostInterval struct {
	SignpostID              int64
	AverageMemory           float64
	CumulativeCPUTime       float64
	CumulativeLogicalWrites float64
	HistogramValues         []HistogramValue
}

type AppSignpost struct {
	Name             string
	Category         string
	MetadataID       int64
	SignpostInterval AppSignpostInterval
}
