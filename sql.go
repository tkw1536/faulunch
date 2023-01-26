package faulunch

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type MenuItem struct {
	ID uint `gorm:"primaryKey"`

	Day      time.Time `gorm:"index"` // the day this item is for
	Location Location  `gorm:"index"` // the location this item is for

	English sql.NullBool `gorm:"not null"` // is this item in english or german?

	Category string // line this item is in
	Title    string // title of this item

	Description string // description of this item
	Beilagen    string // sides

	Preis1 float64 // price (student)
	Preis2 float64 // price (employee)
	Preis3 float64 // price (guest)

	Piktogramme   string // TODO: List of images, properly parsed
	Kj            float64
	Kcal          float64
	Fett          float64
	Gesfett       float64
	Kh            float64
	Zucker        float64
	Ballaststoffe float64
	Eiweiss       float64
	Salz          float64
}

// SyncAll syncs all items from the server
func SyncAll(logger *zerolog.Logger, db *gorm.DB) (failed bool) {
	for _, location := range Locations() {
		if Sync(logger, location, true, db) != nil {
			failed = true
		}
		if Sync(logger, location, false, db) != nil {
			failed = true
		}
	}
	return failed
}

func Sync(logger *zerolog.Logger, location Location, english bool, db *gorm.DB) error {
	// Fetch data from the source
	plan, err := Fetch(location, english)
	logger.Err(err).Str("location", string(location)).Bool("english", english).Msg("fetching data")
	if err != nil {
		return err
	}

	// and sync it!
	return plan.Sync(logger, db, english)
}
