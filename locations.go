package faulunch

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Location string

const (
	CafeteriaBaerenschanzstr Location = "cafeteria-baerenschanzstr"
	CafeteriaBingstr         Location = "cafeteria-bingstr"
	CafeteriaComeIn          Location = "cafeteria-come-in"
	CafeteriaKochstr         Location = "cafeteria-kochstr"
	CafeteriaSuedblick       Location = "cafeteria-suedblick"
	CafeteriaVeilhofstr      Location = "cafeteria-veilhofstr"
	CafeteriaWiso            Location = "cafeteria-wiso"
	MensaAnsbach             Location = "mensa-ansbach"
	MensaEichstaett          Location = "mensa-eichstaett"
	MensaIngolstadt          Location = "mensa-ingolstadt"
	MensaInselschuett        Location = "mensa-inselschuett"
	MensaLmp                 Location = "mensa-lmp"
	MensaRegensburgerstr     Location = "mensa-regensburgerstr"
	MensaSued                Location = "mensa-sued"
	MensateriaOhm            Location = "mensateria-ohm"
	MensateriaStPaul         Location = "mensateria-st-paul"
	MensateriaTriesdorf      Location = "mensateria-triesdorf"
	WohnanlageErwinRommelStr Location = "wohnanlage-erwin-rommel-str"
	WohnanlageHartmannstr    Location = "wohnanlage-hartmannstr"
	WohnanlageStPeter        Location = "wohnanlage-st-peter"
)

func Locations() []Location {
	keys := maps.Keys(locations)
	slices.SortFunc(keys, func(a, b Location) bool {
		return locations[a] < locations[b]
	})
	return keys
}

var locations = map[Location]int{
	MensaSued:                1,
	MensaInselschuett:        2,
	MensaRegensburgerstr:     3,
	MensaAnsbach:             4,
	MensaEichstaett:          5,
	MensateriaOhm:            6,
	MensaIngolstadt:          7,
	MensaLmp:                 8,
	MensateriaStPaul:         9,
	CafeteriaComeIn:          10,
	CafeteriaBaerenschanzstr: 11,
	MensateriaTriesdorf:      12,
	CafeteriaBingstr:         13,
	CafeteriaVeilhofstr:      14,
	CafeteriaKochstr:         17,
	WohnanlageErwinRommelStr: 18,
	WohnanlageHartmannstr:    19,
	CafeteriaSuedblick:       21,
	WohnanlageStPeter:        20,
	CafeteriaWiso:            25,
}

// LocationOfID returns the id of the given location
func LocationOfID(id int) Location {
	for loc, lid := range locations {
		if lid == id {
			return loc
		}
	}
	return ""
}

// IDOfLocation returns the id of a location
func IDOfLocation(l Location) int {
	for loc, id := range locations {
		if loc == l {
			return id
		}
	}
	return 0
}
