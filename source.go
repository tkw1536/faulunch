//spellchecker:words faulunch
package faulunch

//spellchecker:words encoding errors http time github zerolog gorm
import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal/location"
	"github.com/tkw1536/faulunch/internal/plan"
	"gorm.io/gorm"
)

// FetchAndSyncAll fetches and syncs all items into the database.
// It then updates computed fields.
// Returns a boolean indicating failure.
func FetchAndSyncAll(logger *zerolog.Logger, db *gorm.DB) (failed bool) {
	var se SyncEvent
	se.Begin()
	defer func() {
		se.Finish()

		res := se.Store(db)
		logger.Err(res).Msg("logging sync event")
	}()

	for _, loc := range location.Locations() {
		if FetchAndSync(logger, db, loc) != nil {
			failed = true
		}
	}

	if err := RefreshComputedFields(logger, db); err != nil {
		failed = true
	}

	return failed
}

// FetchAndSync is like calling Fetch() and then Sync() for the given location.
func FetchAndSync(logger *zerolog.Logger, db *gorm.DB, loc location.Location) error {
	german, err := plan.Fetch(http.DefaultClient, loc, false)
	logger.Err(err).Str("location", string(loc)).Bool("english", false).Msg("fetching data")
	if err != nil {
		return err
	}

	english, err := plan.Fetch(http.DefaultClient, loc, true)
	logger.Err(err).Str("location", string(loc)).Bool("english", true).Msg("fetching data")
	if err != nil {
		return err
	}

	return Sync(logger, db, german, english)
}

// Sync synchronizes the given german and english plans into the database
// Any previous content for the existing days and locations is erased.
func Sync(logger *zerolog.Logger, db *gorm.DB, german, english plan.Plan) error {
	return db.Transaction(func(tx *gorm.DB) error {
		location, timestamps, items := Merge(logger, german, english)
		times := make([]time.Time, len(timestamps))
		for i, day := range timestamps {
			times[i] = day.Time()
		}

		// delete existing items
		{
			res := tx.
				Where("Day IN ? AND Location = ?", timestamps, location).
				Delete(&MenuItem{})

			logger.Err(res.Error).Int64("count", res.RowsAffected).Str("location", string(location)).Times("timestamps", times).Msg("cleared previous entries")
			if res.Error != nil {
				return res.Error
			}
		}

		if len(items) == 0 {
			logger.Info().Str("location", string(location)).Times("timestamps", times).Msg("no new rows found")
			return nil
		}
		{
			res := tx.Model(&MenuItem{}).Create(&items)
			logger.Err(res.Error).Int64("count", res.RowsAffected).Str("location", string(location)).Times("timestamps", times).Msg("inserted new rows")
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
}

// RefreshComputedFields refreshes all computed fields in the database.
func RefreshComputedFields(logger *zerolog.Logger, db *gorm.DB) error {
	pageSize := 100

	return db.Transaction(func(tx *gorm.DB) error {
		var items []MenuItem

		res := tx.Model(MenuItem{}).FindInBatches(&items, pageSize, func(tx *gorm.DB, batch int) error {
			for i := range items {
				items[i].UpdateComputedFields(logger)
			}

			res := tx.Save(&items)
			logger.Debug().Err(res.Error).Int("batch", batch).Int("count", len(items)).Msg("refreshing computed fields batch")
			return res.Error
		})
		logger.Info().Err(res.Error).Int("rowsAffected", int(res.RowsAffected)).Msg("refreshed computed fields")
		return res.Error
	})
}
