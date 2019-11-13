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
