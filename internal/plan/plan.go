package plan

import (
	"encoding/xml"

	"github.com/tkw1536/faulunch/internal/types"
)

// Plan represents an xml-serializable version of the plan.
// It is directly read from the underlying API.
type Plan struct {
	XMLName  xml.Name `xml:"speiseplan"`
	Location int      `xml:"locationId,attr"`
	Days     []struct {
		Timestamp int `xml:"timestamp,attr"`
		Items     []struct {
			Category string `xml:"category"`
			Title    string `xml:"title"`

			Description string `xml:"description"`
			Beilagen    string `xml:"beilagen"`

			Preis1 types.SmartFloat64 `xml:"preis1"`
			Preis2 types.SmartFloat64 `xml:"preis2"`
			Preis3 types.SmartFloat64 `xml:"preis3"`

			Einheit       string             `xml:"einheit"`
			Piktogramme   string             `xml:"piktogramme"`
			Kj            types.SmartFloat64 `xml:"kj"`
			Kcal          types.SmartFloat64 `xml:"kcal"`
			Fett          types.SmartFloat64 `xml:"fett"`
			Gesfett       types.SmartFloat64 `xml:"gesfett"`
			Kh            types.SmartFloat64 `xml:"kh"`
			Zucker        types.SmartFloat64 `xml:"zucker"`
			Ballaststoffe types.SmartFloat64 `xml:"ballaststoffe"`
			Eiweiss       types.SmartFloat64 `xml:"eiweiss"`
			Salz          types.SmartFloat64 `xml:"salz"`
			Foto          string             `xml:"foto"`
		} `xml:"item"`
	} `xml:"tag"`
}
