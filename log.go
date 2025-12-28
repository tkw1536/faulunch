//spellchecker:words faulunch
package faulunch

//spellchecker:words errors time gorm
import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// SyncEvent represents a synchronization event stored in the database
type SyncEvent struct {
	ID uint `gorm:"primaryKey" json:"-"`

	Start int64 `json:"start"`
	Stop  int64 `json:"stop"`

	Data []byte `json:"-"` // for future usage
}

func (se *SyncEvent) Begin() {
	se.Start = time.Now().Unix()
}

func (se *SyncEvent) Finish() {
	se.Stop = time.Now().Unix()
}

// Store stores this SyncEvent in the database
func (se *SyncEvent) Store(ctx context.Context, db *gorm.DB) error {
	err := gorm.G[SyncEvent](db).Create(ctx, se)
	if err != nil {
		return fmt.Errorf("failed to store sync event: %w", err)
	}
	return nil
}

var ErrNoSync = errors.New("database was never synced")

func (api *API) LastSync(ctx context.Context) (se SyncEvent, err error) {
	return gorm.G[SyncEvent](api.DB).Order("Stop DESC").First(ctx)
}
