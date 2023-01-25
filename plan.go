package faulunch

import (
	"encoding/xml"
	"errors"
	"net/http"
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
		Timestamp string `xml:"timestamp,attr"`
		Item      []struct {
			Category string `xml:"category"`
			Title    string `xml:"title"`

			Description string `xml:"description"`
			Beilagen    string `xml:"beilagen"`

			Preis1 string `xml:"preis1"`
			Preis2 string `xml:"preis2"`
			Preis3 string `xml:"preis3"`

			Einheit       string `xml:"einheit"`
			Piktogramme   string `xml:"piktogramme"`
			Kj            string `xml:"kj"`
			Kcal          string `xml:"kcal"`
			Fett          string `xml:"fett"`
			Gesfett       string `xml:"gesfett"`
			Kh            string `xml:"kh"`
			Zucker        string `xml:"zucker"`
			Ballaststoffe string `xml:"ballaststoffe"`
			Eiweiss       string `xml:"eiweiss"`
			Salz          string `xml:"salz"`
			Foto          string `xml:"foto"`
		} `xml:"item"`
	} `xml:"tag"`
}
