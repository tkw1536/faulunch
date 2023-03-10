package faulunch

import (
	"fmt"
	"html/template"
	"net/url"

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
	keys := maps.Keys(locationIDs)
	slices.SortFunc(keys, func(a, b Location) bool {
		return locationIDs[a] < locationIDs[b]
	})
	return keys
}

func (Location Location) Valid() bool {
	_, ok := locationIDs[Location]
	return ok
}

var locationIDs = map[Location]int{
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

type LocationDescription struct {
	Name string

	Refactory bool
	Cafe      bool
	Internal  bool

	Street   string
	StreetNo string
	ZIP      string
	City     string
}

func (ld LocationDescription) Type(english bool) string {
	if english {
		if ld.Refactory {
			return "Servery"
		}
		if ld.Cafe {
			return "Caf??"
		}
		if ld.Internal {
			return "Internal"
		}
	}

	if ld.Refactory {
		return "Mensa"
	}
	if ld.Cafe {
		return "Caf??"
	}
	if ld.Internal {
		return "Intern"
	}
	return ""
}

func (ld LocationDescription) Address() template.HTML {
	full := fmt.Sprintf("%s %s, %s %s", ld.Street, ld.StreetNo, ld.ZIP, ld.City)
	url := "https://www.openstreetmap.org/search?query=" + url.QueryEscape(full)
	return template.HTML("<a href='" + url + "' rel='noopener noreferer' target='_blank' title='Address'>" + template.HTMLEscapeString(full) + "</a>")
}

func (ld LocationDescription) kind() int {
	if ld.Refactory {
		return 0
	}
	if ld.Cafe {
		return 1
	}
	return 2
}

func (ld LocationDescription) Less(other LocationDescription) bool {

	// kind first
	{
		l, o := ld.kind(), other.kind()
		if l < o {
			return true
		} else if l > o {
			return false
		}
	}

	// city
	{
		l, o := ld.City, other.City
		if l < o {
			return true
		} else if l > o {
			return false
		}
	}

	// title
	{
		l, o := ld.Name, other.Name
		if l < o {
			return true
		} else if l > o {
			return false
		}
	}

	// they are equal
	return true

}

var locationDescriptions = map[Location]LocationDescription{
	MensaSued: {
		Name: "S??dmensa",

		Refactory: true,

		Street:   "Erwin-Rommel-Stra??e",
		StreetNo: "60",
		ZIP:      "91058",
		City:     "Erlangen",
	},
	MensaInselschuett: {
		Name:      "Insel Sch??tt",
		Refactory: true,

		Street:   "Andreij-Sacharow-Platz",
		StreetNo: "1",
		ZIP:      "90403",
		City:     "N??rnberg",
	},
	MensaRegensburgerstr: {
		Name:      "Regensburger Stra??e",
		Refactory: true,

		Street:   "Regensburger Str.",
		StreetNo: "160",
		ZIP:      "90478",
		City:     "N??rnberg",
	},
	MensaAnsbach: {
		Name:      "Ansbach",
		Refactory: true,

		Street:   "Residenzstra??e",
		StreetNo: "8",
		ZIP:      "91522",
		City:     "Ansbach",
	},
	MensaEichstaett: {
		Name:      "Eichst??tt",
		Refactory: true,

		Street:   "Universit??tsallee",
		StreetNo: "2",
		ZIP:      "85072",
		City:     "Eichst??tt",
	},
	MensateriaOhm: {
		Name:      "Mensateria Ohm",
		Refactory: true,

		Street:   "Wollentorstr.",
		StreetNo: "4",
		ZIP:      "90489",
		City:     "N??rnberg",
	},
	MensaIngolstadt: {
		Name:      "Ingolstadt",
		Refactory: true,

		Street:   "Esplanade",
		StreetNo: "10",
		ZIP:      "85049",
		City:     "Ingolstadt",
	},
	MensaLmp: {
		Name:      "Mensa Langemarkplatz",
		Refactory: true,

		Street:   "Langemarckplatz",
		StreetNo: "4",
		ZIP:      "91054",
		City:     "Erlangen",
	},
	MensateriaStPaul: {
		Name:      "Ausgabemensa St. Paul",
		Refactory: true,

		Street:   "Dutzendteichstra??e",
		StreetNo: "24",
		ZIP:      "90478",
		City:     "N??rnberg",
	},
	CafeteriaComeIn: {
		Name: "Cafeteria \"Come IN\" Hohfederstra??e",
		Cafe: true,

		Street:   "Hohfederstra??e",
		StreetNo: "40",
		ZIP:      "90489",
		City:     "N??rnberg",
	},
	CafeteriaBaerenschanzstr: {
		Name: "Cafeteria B??renschanzstra??e",
		Cafe: true,

		Street:   "B??renschanzstr.",
		StreetNo: "4",
		ZIP:      "90429",
		City:     "N??rnberg",
	},
	MensateriaTriesdorf: {
		Name:      "Triesdorf",
		Refactory: true,

		Street:   "Markgrafenstra??e",
		StreetNo: "14",
		ZIP:      "91746",
		City:     "Weidenbach",
	},
	CafeteriaBingstr: {
		Name: "Cafeteria Bingstra??e",
		Cafe: true,

		Street:   "Bingstr.",
		StreetNo: "60",
		ZIP:      "90480",
		City:     "N??rnberg",
	},
	CafeteriaVeilhofstr: {
		Name: "Cafeteria Veilhofstra??e",
		Cafe: true,

		Street:   "Veilhofstra??e",
		StreetNo: "34-40",
		ZIP:      "90489",
		City:     "N??rnberg",
	},
	CafeteriaKochstr: {
		Name: "Cafeteria Kochstra??e",
		Cafe: true,

		Street:   "Kochstr.",
		StreetNo: "4",
		ZIP:      "91054",
		City:     "Erlangen",
	},
	WohnanlageErwinRommelStr: {
		Name:     "Wohnanlage Erwin-Rommel-Stra??e",
		Internal: true,

		Street:   "Erwin-Rommel-Stra??e",
		StreetNo: "51-59",
		ZIP:      "91058",
		City:     "Erlangen",
	},
	WohnanlageHartmannstr: {
		Name:     "Wohnanlage Hartmannstra??e",
		Internal: true,

		Street:   "Hartmannstra??e",
		StreetNo: "125/127/129",
		ZIP:      "91052",
		City:     "Erlangen",
	},
	CafeteriaSuedblick: {
		Name: "Cafeteria S??DBlick",
		Cafe: true,

		Street:   "Erwin-Rommel-Stra??e",
		StreetNo: "51a",
		ZIP:      "91058",
		City:     "Erlangen",
	},
	WohnanlageStPeter: {
		Name:     "Wohnanlage St. Peter",
		Internal: true,

		Street:   "Walter-Meckauer-Stra??e", // Sophienstra??e
		StreetNo: "12-28",                  // 12-16
		ZIP:      "90478",
		City:     "N??rnberg",
	},
	CafeteriaWiso: {
		Name:     "Cafeteria Wiso",
		Internal: true,

		// not sure which one this is
	},
}

// LocationOfID returns the id of the given location
func LocationOfID(id int) Location {
	for loc, lid := range locationIDs {
		if lid == id {
			return loc
		}
	}
	return ""
}

func (l Location) ID() int {
	for loc, id := range locationIDs {
		if loc == l {
			return id
		}
	}
	return 0
}

func (l Location) Description() LocationDescription {
	return locationDescriptions[l]
}
