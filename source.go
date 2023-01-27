package faulunch

import (
	"encoding/xml"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// FetchAndSyncAll fetches and syncs all items into the database.
// Returns a boolean indicating failure.
func FetchAndSyncAll(logger *zerolog.Logger, db *gorm.DB) (failed bool) {
	for _, location := range Locations() {
		if FetchAndSync(logger, db, location) != nil {
			failed = true
		}
	}
	return failed
}

// FetchAndSync is like calling Fetch() and then Sync() for the given location.
func FetchAndSync(logger *zerolog.Logger, db *gorm.DB, location Location) error {
	german, err := Fetch(location, false)
	logger.Err(err).Str("location", string(location)).Bool("english", false).Msg("fetching data")
	if err != nil {
		return err
	}

	english, err := Fetch(location, true)
	logger.Err(err).Str("location", string(location)).Bool("english", true).Msg("fetching data")
	if err != nil {
		return err
	}

	return Sync(logger, db, german, english)
}

var errInvalidStatusCode = errors.New("invalid response code")

// Fetch fetches a plan for the given location and language.
func Fetch(location Location, english bool) (plan Plan, err error) {
	res, err := http.Get(PlanURL(location, english))
	if err != nil {
		return Plan{}, err
	}
	if res.StatusCode != http.StatusOK {
		return Plan{}, errInvalidStatusCode
	}

	err = xml.NewDecoder(res.Body).Decode(&plan)
	return
}

// PlanURL returns the url of a given plan and language
func PlanURL(location Location, english bool) string {
	dest := "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/"
	if english {
		dest += "en/"
	}
	dest += string(location) + ".xml"

	return dest
}

// Sync syncronizes the given german and english plans into the database
// Any previous content for the existing days and locations is erased.
func Sync(logger *zerolog.Logger, db *gorm.DB, german, english Plan) error {
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
