package faulunch

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
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

			Preis1 SmartFloat64 `xml:"preis1"`
			Preis2 SmartFloat64 `xml:"preis2"`
			Preis3 SmartFloat64 `xml:"preis3"`

			Einheit       string       `xml:"einheit"`
			Piktogramme   string       `xml:"piktogramme"`
			Kj            SmartFloat64 `xml:"kj"`
			Kcal          SmartFloat64 `xml:"kcal"`
			Fett          SmartFloat64 `xml:"fett"`
			Gesfett       SmartFloat64 `xml:"gesfett"`
			Kh            SmartFloat64 `xml:"kh"`
			Zucker        SmartFloat64 `xml:"zucker"`
			Ballaststoffe SmartFloat64 `xml:"ballaststoffe"`
			Eiweiss       SmartFloat64 `xml:"eiweiss"`
			Salz          SmartFloat64 `xml:"salz"`
			Foto          string       `xml:"foto"`
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
				menu.Piktogramme.Data = menu.parseIngredients(item.Piktogramme, logger)
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

type SmartFloat64 float64

const (
	smartComma  = ","
	smartPeriod = "."
	smartEmpty  = "-"
)

func (sf64 *SmartFloat64) UnmarshalText(text []byte) error {
	value := string(text)

	// if there are only empty values, ignore them
	if value == "" || value == smartEmpty {
		*sf64 = 0
		return nil
	}

	// replace "," by "." (for german notation)
	if strings.ContainsAny(value, smartComma) && !strings.ContainsAny(value, smartPeriod) {
		value = strings.ReplaceAll(value, smartComma, smartPeriod)
	}

	f64, err := strconv.ParseFloat(value, 64)
	*sf64 = SmartFloat64(f64)
	return err
}

func (sf64 *SmartFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	return sf64.UnmarshalText([]byte(attr.Value))
}
