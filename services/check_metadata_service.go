package services

import (
	"context"

	"github.com/tokopedia/metrickit-backend/models"
	"github.com/jmoiron/sqlx"
)

type CheckMetadataService interface {
	CheckMetadataIfExist(ctx context.Context, metadata *models.AppMetadata) (int64, error)
}

type CheckMetadataServiceImpl struct {
	db *sqlx.DB
}

func NewCheckMetadataService(db *sqlx.DB) *CheckMetadataServiceImpl {
	return &CheckMetadataServiceImpl{
		db: db,
	}
}

func (svc *CheckMetadataServiceImpl) CheckMetadataIfExist(ctx context.Context, metadata *models.AppMetadata) (int64, error) {
	query := `SELECT id FROM metadata WHERE version = ? AND build_number = ? AND device_type = ? AND os = ?`

	stmt, err := svc.db.PreparexContext(ctx, query)

	if err != nil {
		return 0, err
	}

	rows, err := stmt.QueryxContext(ctx, metadata.Version, metadata.BuildNumber, metadata.DeviceType, metadata.OSVersion)

	defer rows.Close()

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)

		if err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, nil // if no rows returned
}
