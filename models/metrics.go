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

type AppSignpostHistogram struct {
	SignpostIntervalID int64
	HistogramValues    []HistogramValue
}

type AppSignpostInterval struct {
	SignpostID              int64
	AverageMemory           float64
	CumulativeCPUTime       float64
	CumulativeLogicalWrites float64
	SignpostHistogram       AppSignpostHistogram
}

type AppSignpost struct {
	Name             string
	Category         string
	MetadataID       int64
	SignpostInterval AppSignpostInterval
}

type AppHangTime struct {
	HistogramValue
	MetadataID int64
}
