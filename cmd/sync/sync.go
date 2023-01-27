package main

import (
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: cmd/sync <path-to-db>")
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	log := zerolog.New(output).With().Timestamp().Logger()

	// open the database
	db, err := gorm.Open(sqlite.Open(os.Args[1]), &gorm.Config{})
	log.Err(err).Msg("opening database")
	if err != nil {
		panic(err)
	}

	// register a close once we're done
	{
		db, err := db.DB()
		log.Err(err).Msg("registering shutdown")
		if err != nil {
			panic(err)
		}
		defer db.Close()
	}

	// do the migration
	{
		err := db.AutoMigrate(&faulunch.MenuItem{})
		log.Err(err).Msg("migrating database")
		if err != nil {
			panic(err)
		}
	}

	// fetch all the items
	{
		failed := faulunch.FetchAndSyncAll(&log, db)
		if failed {
			panic("failed to sync all locations")
		}
	}

}
