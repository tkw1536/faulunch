//spellchecker:words faulunch
package faulunch

//spellchecker:words errors slices gorm
import (
	"errors"

	"slices"

	"github.com/tkw1536/faulunch/internal/ltime"
	"gorm.io/gorm"
)

type API struct {
	DB *gorm.DB
}

// Locations returns the list of available locations in the database.
// They are sorted by their type.
func (api *API) Locations() (locations []Location, err error) {
	res := api.DB.Model(&MenuItem{}).Distinct("Location").Pluck("location", &locations)
	err = res.Error

	slices.SortStableFunc(locations, func(a, b Location) int {
		return a.Description().Cmp(b.Description())
	})

	return
}

// Pagination represents an environment for the
type Pagination struct {
	Head []ltime.Day

	Prev []ltime.Day
	Now  ltime.Day
	Next []ltime.Day

	Tail []ltime.Day
}

var errNoCurrent = errors.New("no current item")

func (api *API) CurrentDay(location Location, query ltime.Day) (day ltime.Day, err error) {
	var days []ltime.Day
	res := api.DB.Model(&MenuItem{}).
		Where("Location = ? AND Day <= ? ", location, query).
		Distinct("Day").Order("day DESC").
		Limit(1).
		Pluck("day", &days)
	if len(days) == 0 {
		return day, errNoCurrent
	}

	if res.Error != nil {
		return day, res.Error
	}

	return days[0], nil
}

// DayPagination builds up day-based pagination for the given query.
func (api *API) DayPagination(location Location, query ltime.Day, size int) (pagination Pagination, err error) {
	if size < 1 {
		panic("DayPagination: limit < 1")
	}

	// make a query to fill everything up
	for _, q := range []*gorm.DB{
		api.DB.Model(&MenuItem{}).Where("Location = ? AND Day <= ? ", location, query).
			Distinct("Day").Order("day DESC").
			Limit(size+1).
			Pluck("day", &pagination.Prev),
		api.DB.Model(&MenuItem{}).
			Where("Location = ? AND Day > ? ", location, query).
			Distinct("Day").Order("day ASC").
			Limit(size).
			Pluck("day", &pagination.Next),
		api.DB.Model(&MenuItem{}).
			Where("Location = ?", location).
			Distinct("Day").Order("day ASC").
			Limit(size).
			Pluck("day", &pagination.Head),
		api.DB.Model(&MenuItem{}).
			Where("Location = ?", location).
			Distinct("Day").Order("day DESC").
			Limit(size).
			Pluck("day", &pagination.Tail),
	} {
		if q.Error != nil {
			return pagination, q.Error
		}
	}

	// errNoCurrent if there is no current item
	if len(pagination.Prev) == 0 {
		return pagination, errNoCurrent
	}

	// pick the current item as the last of the previous items
	// which is *first* in the array because it is DSC.
	pagination.Now = pagination.Prev[0]
	pagination.Prev = pagination.Prev[1:]

	// reverse everything into the right order
	reverse(pagination.Prev)
	reverse(pagination.Tail)

	// remove invalid values from head
	{

		// determine the limit for values to be the first value in prev
		// or (if that is empty) the current day
		limit := pagination.Now
		if len(pagination.Prev) > 0 {
			limit = pagination.Prev[0]
		}

		// filter the head value to only contain value smaller than the limit

		for i, d := range pagination.Head {
			if d >= limit {
				// we abuse the fact that it is assumed to be sorted
				// so we know all further elements are >= also
				pagination.Head = pagination.Head[:i]
				break
			}
		}
	}

	// remove duplicate values from tail
	{

		// determine the limit for values to be the last value in .Next
		// or (if that is empty) the current day
		limit := pagination.Now
		if len(pagination.Next) > 0 {
			limit = pagination.Next[len(pagination.Next)-1]
		}

		// find the first valid value in the tail
		valid := false
		for i, d := range pagination.Tail {
			if d > limit {
				// and take only values from then onwards
				pagination.Tail = pagination.Tail[i:]
				valid = true
				break
			}
		}

		// no element in the tail was actually valid
		// (everything already contained in .Next)
		if !valid {
			pagination.Tail = nil
		}
	}

	return
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// KnowsLocation checks if at least one time for the given location is known
func (api *API) KnowsLocation(location Location) (exists bool, err error) {
	err = api.DB.Model(&MenuItem{}).Where("Location = ?", location).Select("count(*) > 0").Find(&exists).Error
	return
}

// Days returns the days with an explicit menu starting at day, and including (at most) the next count days.
func (api *API) Days(location Location, day ltime.Day, count int) (days []ltime.Day, err error) {
	start := day.Normalize()
	end := day.Add(count)

	res := api.DB.Model(&MenuItem{}).Where("Location = ? AND day >= ? AND day < ?", location, start, end).Order("day DESC").Distinct().Pluck("day", &days)
	err = res.Error
	return
}

// MenuItems returns the menu items for the given day and time.
// They are sorted by category.
// If it does not exist, an empty menu item is returned.
func (api *API) MenuItems(location Location, day ltime.Day) (items []MenuItem, err error) {
	res := api.DB.Model(&MenuItem{}).Where("Location = ? AND day = ?", location, day).Order("Category ASC").Find(&items)
	slices.SortStableFunc(items, func(a, b MenuItem) int { return a.Cmp(b) })
	err = res.Error
	return
}
