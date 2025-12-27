package plan

import (
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/tkw1536/faulunch/internal/location"
)

var errInvalidStatusCode = errors.New("invalid response code")

// Fetch fetches a plan for the given location and language.
func Fetch(loc location.Location, english bool) (plan Plan, err error) {
	res, err := http.Get(PlanURL(loc, english))
	if err != nil {
		return plan, err
	}
	if res.StatusCode != http.StatusOK {
		return plan, errInvalidStatusCode
	}

	err = xml.NewDecoder(res.Body).Decode(&plan)
	return
}

// PlanURL returns the url of a given plan and language
func PlanURL(loc location.Location, english bool) string {
	dest := "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/"
	if english {
		dest += "en/"
	}
	dest += string(loc) + ".xml"

	return dest
}
