package plan

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/tkw1536/faulunch/internal/location"
)

var errInvalidStatusCode = errors.New("invalid response code")

// ClientLike is an interface that is implemented by [*http.Client].
type ClientLike interface {
	// Do sends an HTTP request and returns an HTTP response, following
	// policy (such as redirects, cookies, auth) as configured on the
	// client.
	//
	// See [net/http.Client.Do] for more details.
	Do(req *http.Request) (*http.Response, error)
}

// Fetch fetches a plan for the given location and language.
func Fetch(ctx context.Context, client ClientLike, loc location.Location, english bool) (plan Plan, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, PlanURL(loc, english), nil)
	if err != nil {
		return plan, fmt.Errorf("failed to create request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return plan, fmt.Errorf("failed to fetch plan: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return plan, fmt.Errorf("failed to fetch plan: %w", errInvalidStatusCode)
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
