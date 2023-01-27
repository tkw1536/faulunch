package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch"
	"gorm.io/gorm"
)

func main() {
	args := flag.Args()
	if len(args) != 1 {
		panic("Usage: cmd/serve [...flags] <path-to-db>")
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	log := zerolog.New(output).With().Timestamp().Logger()

	if flagDebug {
		log = log.Level(zerolog.DebugLevel)
	} else {
		log = log.Level(zerolog.InfoLevel)
	}

	// open the database
	db, err := gorm.Open(sqlite.Open(args[0]), &gorm.Config{})
	log.Err(err).Msg("opening database")
	if err != nil {
		panic(err)
	}

	// start automatically syncing if requested
	if flagAutoSync > 0 {

		// do the migration
		{
			err := db.AutoMigrate(&faulunch.MenuItem{})
			log.Err(err).Msg("migrating database")
			if err != nil {
				panic(err)
			}
		}

		go func() {
			for {
				failed := faulunch.FetchAndSyncAll(&log, db)
				if failed {
					log.Error().Msg("failed to sync")
				}
				time.Sleep(flagAutoSync)
			}
		}()
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

	// create an api and start serving it
	{
		server := faulunch.Server{
			API: &faulunch.API{
				DB: db,
			},
			Logger: &log,
			Legal: faulunch.ServerLegal{
				Link:     flagLink,
				DEString: flagDEText,
				ENString: flagENText,
			},
		}

		log.Info().Str("addr", flagAddr).Msg("server listening")
		err := http.ListenAndServe(flagAddr, &server)
		log.Err(err).Str("addr", flagAddr).Msg("server failed to listen")
	}

}

var flagAutoSync time.Duration
var flagDebug bool = false
var flagAddr string = "127.0.0.1:3000"
var flagLink string = "https://privacy.kwarc.info/"
var flagDEText string = "Keine offizielle Seite des Studentenwerks. Alle Angaben, insbesondere zu Speiseplänen und Preisen, sind ohne Gewähr. Siehe auch Impressum und Datenschutz. "
var flagENText string = "Not an official page of Studentenwerk. All information subject to change. See also Imprint and Privacy Policy. "

func init() {
	defer flag.Parse()

	flag.DurationVar(&flagAutoSync, "sync", flagAutoSync, "automatically sync")
	flag.BoolVar(&flagDebug, "debug", flagDebug, "Set debug log level")

	flag.StringVar(&flagAddr, "addr", flagAddr, "Address to bind to")
	flag.StringVar(&flagDEText, "legal-de", flagDEText, "text for german legal link")
	flag.StringVar(&flagDEText, "legal-en", flagENText, "text for english legal link")
	flag.StringVar(&flagLink, "legal-link", flagLink, "url for legal link")
}
