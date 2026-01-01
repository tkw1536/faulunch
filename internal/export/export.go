package export

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// NewExporter creates a new sqlite exporter.
//
// db is a database handle pointing to an sqlite database.
// query is a query (to be run on the database) that returns a unique identifier for the current contents of the database.
// if this changes, any changes are automatically invalidated.
//
// The returned function when called will write a consistent copy of the sqlite database to the given writer.
// The database is internally created using the VACUUM INTO command.
//
// The returned function may be called multiple times, and will cache the created result on disk in a temporary file.
// When the context is cancelled, or the close function is called, the temporary file is deleted.
// The close function waits for the database to be deleted.
//
// The returned function may be called concurrently.
func NewExporter(ctx context.Context, logger *zerolog.Logger, db *sql.DB, query string) (write func(w http.ResponseWriter, r *http.Request) error, close func() error) {
	e := &exporter{
		ctx:    ctx,
		logger: logger,
		db:     db,
		query:  query,
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-ctx.Done()
		e.cleanup()
	}()

	return e.write, func() error {
		cancel()
		e.cleanup()
		return nil
	}
}

type exporter struct {
	ctx    context.Context
	logger *zerolog.Logger

	db    *sql.DB
	query string

	mu         sync.RWMutex
	lastID     string
	tempFile   string
	exportTime time.Time
}

// cleanup removes the temporary file if it exists
func (e *exporter) cleanup() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.tempFile == "" {
		return nil
	}

	e.logger.Debug().Str("path", e.tempFile).Msg("deleting export file")
	if err := os.Remove(e.tempFile); err != nil {
		return fmt.Errorf("failed to remove temporary file: %w", err)
	}

	e.tempFile = ""
	e.lastID = ""
	return nil
}

// getCurrentID runs the query and returns the current identifier
func (e *exporter) getCurrentID() (string, error) {
	var id string
	err := e.db.QueryRowContext(e.ctx, e.query).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// ensureExport ensures that we have an up-to-date export file.
// Returns the path to the temporary file and the time when the export was created.
func (e *exporter) ensureExport() (string, time.Time, error) {
	currentID, err := e.getCurrentID()
	if err != nil {
		return "", time.Time{}, err
	}

	// Fast path: check if cached version is still valid
	e.mu.RLock()
	if e.tempFile != "" && e.lastID == currentID {
		path := e.tempFile
		exportTime := e.exportTime
		e.mu.RUnlock()
		e.logger.Debug().Str("path", path).Msg("reusing export file")
		return path, exportTime, nil
	}
	e.mu.RUnlock()

	// Slow path: need to create new export
	e.mu.Lock()
	defer e.mu.Unlock()

	// Double-check after acquiring write lock
	if e.tempFile != "" && e.lastID == currentID {
		e.logger.Debug().Str("path", e.tempFile).Msg("reusing export file")
		return e.tempFile, e.exportTime, nil
	}

	// Remove old temp file if exists
	if e.tempFile != "" {
		e.logger.Debug().Str("path", e.tempFile).Msg("deleting old export file")
		if err := os.Remove(e.tempFile); err != nil {
			return "", time.Time{}, fmt.Errorf("failed to remove old temporary file: %w", err)
		}
		e.tempFile = ""
		e.lastID = ""
		e.exportTime = time.Time{}
	}

	// Create new temp file
	tmpFile, err := os.CreateTemp("", "sqlite-export-*.db")
	if err != nil {
		return "", time.Time{}, err
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Remove it first, VACUUM INTO requires the file to not exist
	if err := os.Remove(tmpPath); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to remove temporary file path before VACUUM INTO: %w", err)
	}

	// Run VACUUM INTO
	_, err = e.db.ExecContext(e.ctx, "VACUUM INTO ?", tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return "", time.Time{}, err
	}

	e.tempFile = tmpPath
	e.lastID = currentID
	e.exportTime = time.Now()
	e.logger.Debug().Str("path", tmpPath).Msg("created export file")
	return tmpPath, e.exportTime, nil
}

// write writes the database export to the given http.ResponseWriter using http.ServeContent
func (e *exporter) write(w http.ResponseWriter, r *http.Request) error {
	if e.ctx.Err() != nil {
		return e.ctx.Err()
	}

	path, exportTime, err := e.ensureExport()
	if err != nil {
		return err
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	// Check that path is still valid
	if e.tempFile != path {
		return errors.New("export invalidated during write")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open export file: %w", err)
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/x-sqlite3")
	http.ServeContent(w, r, "export.db", exportTime, f)
	return nil
}
