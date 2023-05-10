package faulunch

import (
	"errors"
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
func (se *SyncEvent) Store(db *gorm.DB) error {
	res := db.Model(&SyncEvent{}).Create(se)
	return res.Error
}

var ErrNoSync = errors.New("database was never synced")

func (api *API) LastSync() (se SyncEvent, err error) {
	res := api.DB.Model(&SyncEvent{}).Order("Stop DESC").First(&se)
	if res.Error != nil {
		return se, res.Error
	}

	if res.RowsAffected == 0 {
		return se, ErrNoSync
	}

	return se, nil
}
