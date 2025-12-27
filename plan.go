//spellchecker:words faulunch
package faulunch

//spellchecker:words encoding strconv strings github zerolog faulunch internal
import (
	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal"
	"github.com/tkw1536/faulunch/internal/ltime"
	"github.com/tkw1536/faulunch/internal/plan"
	"github.com/tkw1536/faulunch/internal/types"
)

// Merge merges the german and english plan for the given location into a set of menu items, a location, and a timestamp.
func Merge(logger *zerolog.Logger, german plan.Plan, english plan.Plan) (location Location, timestamps []ltime.Day, menu []MenuItem) {
	// extract the location
	location = LocationOfID(german.Location)

	// create a map from days => category to menuitem
	dayCatMap := make(map[int]map[string]MenuItem)

	for _, daylang := range []struct {
		Plan    plan.Plan
		English bool
	}{
		{german, false},
		{english, true},
	} {
		for _, day := range daylang.Plan.Days {
			timestamp := ltime.ParseDay(day.Timestamp)

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

				menu.Preis1 = types.LPrice(item.Preis1)
				menu.Preis2 = types.LPrice(item.Preis2)
				menu.Preis3 = types.LPrice(item.Preis3)

				// TODO: Extract Piktogramme
				internal.SetJSONData(&menu.Piktogramme, menu.parseIngredients(item.Piktogramme, logger))
				menu.Kj = types.LFloat(item.Kj)
				menu.Kcal = types.LFloat(item.Kcal)
				menu.Fett = types.LFloat(item.Fett)
				menu.Gesfett = types.LFloat(item.Gesfett)
				menu.Kh = types.LFloat(item.Kh)
				menu.Zucker = types.LFloat(item.Zucker)
				menu.Ballaststoffe = types.LFloat(item.Ballaststoffe)
				menu.Eiweiss = types.LFloat(item.Eiweiss)
				menu.Salz = types.LFloat(item.Salz)

				catMap[item.Category] = menu
			}

			dayCatMap[day.Timestamp] = catMap
		}
	}

	// build the menu and all the timestamps
	timestamps = make([]ltime.Day, 0, len(dayCatMap))
	for t, catMap := range dayCatMap {
		timestamps = append(timestamps, ltime.ParseDay(t))
		for _, mitem := range catMap {
			menu = append(menu, mitem)
		}
	}

	return
}
