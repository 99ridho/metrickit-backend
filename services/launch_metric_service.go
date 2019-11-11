package services

import (
	"context"

	"github.com/99ridho/metrickit-backend/models"
	"github.com/jmoiron/sqlx"
)

type LaunchMetricService interface {
	Store(ctx context.Context, launchTimes []*models.AppLaunchTime, metadata *models.AppMetadata) ([]int64, error)
}

type LaunchMetricServiceImpl struct {
	db *sqlx.DB
}

func NewLaunchMetricService(db *sqlx.DB) *LaunchMetricServiceImpl {
	return &LaunchMetricServiceImpl{
		db: db,
	}
}

func (svc *LaunchMetricServiceImpl) Store(ctx context.Context, launchTimes []*models.AppLaunchTime, metadata *models.AppMetadata) ([]int64, error) {
	metadataID, err := svc.checkMetadataIfExist(ctx, metadata)

	if err != nil {
		return []int64{}, err
	}

	query := `INSERT INTO app_launch_time_first_draw (metadata_id, range_start, range_end, frequency) VALUES (?,?,?,?)`

	if metadataID != 0 {
		ids := make([]int64, 0)

		for _, launchTime := range launchTimes {
			launchTime.MetadataID = metadataID

			stmt, err := svc.db.PreparexContext(ctx, query)

			if err != nil {
				return []int64{}, err
			}

			result, err := stmt.ExecContext(ctx, launchTime.MetadataID, launchTime.RangeStart, launchTime.RangeEnd, launchTime.Frequency)

			if err != nil {
				return []int64{}, err
			}

			lastID, _ := result.LastInsertId()

			ids = append(ids, lastID)
		}

		return ids, nil
	}

	// if no metadata inserted, do transaction to insert both metadata and respective metrics
	tx, err := svc.db.BeginTxx(ctx, nil)

	if err != nil {
		tx.Rollback()
		return []int64{}, err
	}

	// insert metadata
	metadataInsertQuery := `INSERT INTO metadata (version, build_number, device_type, os) VALUES (?,?,?,?)`

	stmt, err := tx.PreparexContext(ctx, metadataInsertQuery)

	if err != nil {
		tx.Rollback()
		return []int64{}, err
	}

	result, err := stmt.ExecContext(ctx, metadata.Version, metadata.BuildNumber, metadata.DeviceType, metadata.OSVersion)

	if err != nil {
		tx.Rollback()
		return []int64{}, err
	}

	lastMetadataInsertID, _ := result.LastInsertId()

	ids := make([]int64, 0)
	for _, launchTime := range launchTimes {
		launchTime.MetadataID = lastMetadataInsertID

		stmt, err := tx.PreparexContext(ctx, query)

		if err != nil {
			tx.Rollback()
			return []int64{}, err
		}

		result, err := stmt.ExecContext(ctx, launchTime.MetadataID, launchTime.RangeStart, launchTime.RangeEnd, launchTime.Frequency)

		if err != nil {
			tx.Rollback()
			return []int64{}, err
		}

		lastID, _ := result.LastInsertId()

		ids = append(ids, lastID)
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return []int64{}, err
	}

	return ids, nil
}

func (svc *LaunchMetricServiceImpl) checkMetadataIfExist(ctx context.Context, metadata *models.AppMetadata) (int64, error) {
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
