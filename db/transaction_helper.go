package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionHelper struct {
	DB *sqlx.DB
}

func (th *TransactionHelper) DoTransaction(ctx context.Context, handler func(tx *sqlx.Tx) error) (err error) {
	tx, err := th.DB.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			return
		}
	}()

	err = handler(tx)
	return
}
