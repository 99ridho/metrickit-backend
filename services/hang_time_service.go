package services

import (
	"context"

	database "github.com/99ridho/metrickit-backend/db"
	"github.com/99ridho/metrickit-backend/models"
	"github.com/jmoiron/sqlx"
)

type HangTimeService interface {
	Store(ctx context.Context, hangRates []*models.AppHangTime, metadata *models.AppMetadata) ([]int64, error)
}

type HangTimeServiceImpl struct {
	db                   *sqlx.DB
	transactionHelper    *database.TransactionHelper
	checkMetadataService CheckMetadataService
}

func NewHangTimeService(db *sqlx.DB) *HangTimeServiceImpl {
	return &HangTimeServiceImpl{
		db: db,
		transactionHelper: &database.TransactionHelper{
			DB: db,
		},
		checkMetadataService: NewCheckMetadataService(db),
	}
}

func (svc *HangTimeServiceImpl) Store(ctx context.Context, hangRates []*models.AppHangTime, metadata *models.AppMetadata) ([]int64, error) {
	existMetadataID, err := svc.checkMetadataService.CheckMetadataIfExist(ctx, metadata)

	if err != nil {
		return []int64{}, err
	}

	ids := make([]int64, 0)

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

		query := "INSERT INTO app_hang_time (metadata_id, range_start, range_end, frequency) VALUES (?,?,?,?)"
		for _, ht := range hangRates {
			ht.MetadataID = metadataID

			stmt, err := tx.PreparexContext(ctx, query)

			if err != nil {
				return err
			}

			result, err := stmt.ExecContext(ctx, ht.MetadataID, ht.RangeStart, ht.RangeEnd, ht.Frequency)

			if err != nil {
				return err
			}

			lastID, _ := result.LastInsertId()

			ids = append(ids, lastID)
		}

		return nil
	})

	if err != nil {
		return []int64{}, err
	}

	return ids, nil
}
