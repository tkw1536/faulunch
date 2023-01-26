package faulunch

import (
	"database/sql"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (plan Plan) Sync(logger *zerolog.Logger, db *gorm.DB, english bool) error {
	return db.Transaction(func(tx *gorm.DB) error {
		location, timestamps, items := plan.Items(english)

		// delete existing items
		{
			res := tx.
				Where("Day IN ? AND Location = ? AND English = ?", timestamps, location, sql.NullBool{Valid: true, Bool: english}).
				Delete(&MenuItem{})

			logger.Err(res.Error).Int64("count", res.RowsAffected).Str("location", string(location)).Times("timestamps", timestamps).Msg("cleared previous entries")
			if res.Error != nil {
				return res.Error
			}
		}

		if len(items) == 0 {
			logger.Info().Str("location", string(location)).Times("timestamps", timestamps).Msg("no new rows found")
			return nil
		}
		{
			res := tx.Model(&MenuItem{}).Create(&items)
			logger.Err(res.Error).Int64("count", res.RowsAffected).Str("location", string(location)).Times("timestamps", timestamps).Msg("inserted new rows")
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})

}
