package models

type AppLaunchTime struct {
	ID         int64   `db:"id" json:"id"`
	MetadataID int64   `db:"metadata_id" json:"metadata_id"`
	RangeStart float64 `db:"range_start" json:"range_start"`
	RangeEnd   float64 `db:"range_end" json:"range_end"`
	Frequency  int64   `db:"frequency" json:"frequency"`
}
