package services

import (
	"context"

	database "github.com/tokopedia/metrickit-backend/db"
	"github.com/tokopedia/metrickit-backend/models"
	"github.com/jmoiron/sqlx"
)

type SignpostService interface {
	Store(ctx context.Context, signposts []*models.AppSignpost, metadata *models.AppMetadata) ([]int64, error)
}

type SignpostServiceImpl struct {
	db                   *sqlx.DB
	transactionHelper    *database.TransactionHelper
	checkMetadataService CheckMetadataService
}

func NewSignpostService(db *sqlx.DB) *SignpostServiceImpl {
	return &SignpostServiceImpl{
		db: db,
		transactionHelper: &database.TransactionHelper{
			DB: db,
		},
		checkMetadataService: NewCheckMetadataService(db),
	}
}

func (svc *SignpostServiceImpl) Store(ctx context.Context, signposts []*models.AppSignpost, metadata *models.AppMetadata) ([]int64, error) {
	existMetadataID, err := svc.checkMetadataService.CheckMetadataIfExist(ctx, metadata)

	if err != nil {
		return []int64{}, err
	}

	insertedSignpostIDs := make([]int64, 0)

	err = svc.transactionHelper.DoTransaction(ctx, func(tx *sqlx.Tx) error {
		var metadataID int64
		if existMetadataID != 0 {
			metadataID = existMetadataID
		} else {
			metadataInsertQuery := `INSERT INTO metadata (version, build_number, device_type, os) VALUES (?,?,?,?)`

			stmt, err := tx.PreparexContext(ctx, metadataInsertQuery)

			if err != nil {
				return err
			}

			result, err := stmt.ExecContext(ctx, metadata.Version, metadata.BuildNumber, metadata.DeviceType, metadata.OSVersion)

			if err != nil {
				return err
			}

			metadataID, _ = result.LastInsertId()
		}

		for _, signpost := range signposts {
			signpost.MetadataID = metadataID

			signpostID, err := svc.insertSignpost(ctx, tx, signpost)

			if err != nil {
				return err
			}

			insertedSignpostIDs = append(insertedSignpostIDs, signpostID)
		}

		return nil
	})

	if err != nil {
		return []int64{}, err
	}

	return insertedSignpostIDs, nil
}

func (svc *SignpostServiceImpl) insertSignpost(ctx context.Context, tx *sqlx.Tx, signpost *models.AppSignpost) (int64, error) {
	existingSignpostID, err := svc.checkSignpostIfExist(ctx, signpost)

	if err != nil {
		return 0, err
	}

	var signpostID int64
	if existingSignpostID != 0 {
		signpostID = existingSignpostID
	} else {
		signpostInsertQuery := "INSERT INTO app_signpost (metadata_id, name, category) VALUES (?,?,?)"

		stmt, err := tx.PreparexContext(ctx, signpostInsertQuery)

		if err != nil {
			return 0, err
		}

		signpostInsertResult, err := stmt.ExecContext(ctx, signpost.MetadataID, signpost.Name, signpost.Category)

		if err != nil {
			return 0, err
		}

		signpostID, _ = signpostInsertResult.LastInsertId()
	}

	signpost.SignpostInterval.SignpostID = signpostID

	err = svc.insertSignpostInterval(ctx, tx, signpost)

	if err != nil {
		return 0, err
	}

	return signpostID, nil
}

func (svc *SignpostServiceImpl) checkSignpostIfExist(ctx context.Context, signpost *models.AppSignpost) (int64, error) {
	query := "SELECT id FROM app_signpost WHERE metadata_id = ? AND name = ? AND category = ?"

	stmt, err := svc.db.PreparexContext(ctx, query)

	if err != nil {
		return 0, err
	}

	rows, err := stmt.QueryxContext(ctx, signpost.MetadataID, signpost.Name, signpost.Category)

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

	return 0, nil
}

func (svc *SignpostServiceImpl) insertSignpostInterval(ctx context.Context, tx *sqlx.Tx, signpost *models.AppSignpost) error {
	signpostIntervalQuery := "INSERT INTO app_signpost_interval (signpost_id, average_memory, cumulative_cpu_time, cumulative_logical_writes) VALUES (?,?,?,?)"

	stmt, err := tx.PreparexContext(ctx, signpostIntervalQuery)

	if err != nil {
		return err
	}

	signpostIntervalResult, err := stmt.ExecContext(
		ctx,
		signpost.SignpostInterval.SignpostID,
		signpost.SignpostInterval.AverageMemory,
		signpost.SignpostInterval.CumulativeCPUTime,
		signpost.SignpostInterval.CumulativeLogicalWrites,
	)

	if err != nil {
		return err
	}

	signpostIntervalID, _ := signpostIntervalResult.LastInsertId()

	signpost.SignpostInterval.SignpostHistogram.SignpostIntervalID = signpostIntervalID

	for _, histogramValue := range signpost.SignpostInterval.SignpostHistogram.HistogramValues {
		histogramInsertQuery := "INSERT INTO app_signpost_histogram (signpost_interval_id, range_start, range_end, frequency) VALUES (?,?,?,?)"

		stmt, err := tx.PreparexContext(ctx, histogramInsertQuery)

		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(
			ctx,
			signpost.SignpostInterval.SignpostHistogram.SignpostIntervalID,
			histogramValue.RangeStart,
			histogramValue.RangeEnd,
			histogramValue.Frequency,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
