package internal

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func SetJSONData[T any](elem *datatypes.JSONType[T], value T) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return elem.UnmarshalJSON(bytes)
}
