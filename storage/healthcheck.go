package storage

import (
	"context"
)

func (db *wrappedDB) Check(ctx context.Context) error {
	rows, err := db.QueryContext(ctx, "SELECT 1")
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
