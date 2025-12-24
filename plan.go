//spellchecker:words faulunch
package faulunch

//spellchecker:words encoding strconv strings github zerolog faulunch internal
import (
	"encoding/xml"

	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal"
	"github.com/tkw1536/faulunch/internal/types"
)

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

// Merge merges the german and english plan for the given location into a set of menu items, a location, and a timestamp.
func Merge(logger *zerolog.Logger, german Plan, english Plan) (location Location, timestamps []Day, menu []MenuItem) {
	// extract the location
	location = LocationOfID(german.Location)

	// create a map from days => category to menuitem
	dayCatMap := make(map[int]map[string]MenuItem)

	for _, daylang := range []struct {
		Plan    Plan
		English bool
	}{
		{german, false},
		{english, true},
	} {
		for _, day := range daylang.Plan.Days {
			timestamp := ParseDay(day.Timestamp)

			// generate a map of categories
			catMap := dayCatMap[day.Timestamp]
			if catMap == nil {
				catMap = make(map[string]MenuItem, len(day.Items))
			}

			for _, item := range day.Items {
				menu := catMap[item.Category]
				menu.Location = location
				menu.Day = timestamp
				menu.Category = item.Category

				if daylang.English {
					menu.TitleEN = item.Title
					menu.DescriptionEN = item.Description
					menu.BeilagenEN = item.Beilagen
				} else {
					menu.TitleDE = item.Title
					menu.DescriptionDE = item.Description
					menu.BeilagenDE = item.Beilagen
				}

				menu.Preis1 = LPrice(item.Preis1)
				menu.Preis2 = LPrice(item.Preis2)
				menu.Preis3 = LPrice(item.Preis3)

				// TODO: Extract Piktogramme
				internal.SetJSONData(&menu.Piktogramme, menu.parseIngredients(item.Piktogramme, logger))
				menu.Kj = LFloat(item.Kj)
				menu.Kcal = LFloat(item.Kcal)
				menu.Fett = LFloat(item.Fett)
				menu.Gesfett = LFloat(item.Gesfett)
				menu.Kh = LFloat(item.Kh)
				menu.Zucker = LFloat(item.Zucker)
				menu.Ballaststoffe = LFloat(item.Ballaststoffe)
				menu.Eiweiss = LFloat(item.Eiweiss)
				menu.Salz = LFloat(item.Salz)

				catMap[item.Category] = menu
			}

			dayCatMap[day.Timestamp] = catMap
		}
	}

	// build the menu and all the timestamps
	timestamps = make([]Day, 0, len(dayCatMap))
	for t, catMap := range dayCatMap {
		timestamps = append(timestamps, ParseDay(t))
		for _, mitem := range catMap {
			menu = append(menu, mitem)
		}
	}

	return
}
