package models

type AppMetadata struct {
	ID          int64  `db:"id" json:"id"`
	Version     string `db:"version" json:"version"`
	BuildNumber string `db:"build_number" json:"build_number"`
	DeviceType  string `db:"device_type" json:"device_type"`
	OSVersion   string `db:"os" json:"os"`
}
