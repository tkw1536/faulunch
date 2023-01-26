package faulunch

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PlanURL returns the plan with the given url
func PlanURL(location Location, english bool) string {
	dest := "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/"
	if english {
		dest += "en/"
	}
	dest += string(location) + ".xml"

	return dest
}

var errInvalidStatusCode = errors.New("invalid response code")

func Fetch(location Location, english bool) (plan Plan, err error) {
	res, err := http.Get(PlanURL(location, english))
	if err != nil {
		return Plan{}, err
	}
	if res.StatusCode != http.StatusOK {
		return Plan{}, errInvalidStatusCode
	}

	err = xml.NewDecoder(res.Body).Decode(&plan)
	return
}

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

func (plan Plan) Items(english bool) (location Location, timestamps []time.Time, menu []MenuItem) {
	location = LocationOfID(plan.Location)
	ebool := sql.NullBool{Bool: english, Valid: true}
	for _, day := range plan.Days {
		timestamp := time.Unix(int64(day.Timestamp), 0)
		timestamps = append(timestamps, timestamp)

		for _, item := range day.Items {
			mitem := MenuItem{
				Day:      timestamp,
				Location: location,

				English: ebool,

				Category: item.Category,
				Title:    item.Title,

				Description: item.Description,
				Beilagen:    item.Beilagen,

				Preis1: float64(item.Preis1),
				Preis2: float64(item.Preis2),
				Preis3: float64(item.Preis3),

				Piktogramme:   item.Piktogramme,
				Kj:            float64(item.Kj),
				Kcal:          float64(item.Kcal),
				Fett:          float64(item.Fett),
				Gesfett:       float64(item.Gesfett),
				Kh:            float64(item.Kh),
				Zucker:        float64(item.Zucker),
				Ballaststoffe: float64(item.Ballaststoffe),
				Eiweiss:       float64(item.Eiweiss),
				Salz:          float64(item.Salz),
			}

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
