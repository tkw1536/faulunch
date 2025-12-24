package annotations

import (
	"html/template"

	"github.com/tkw1536/faulunch/internal/fmap"
)

type Additive string

const (
	Color           Additive = "1"
	Caffeine        Additive = "2"
	Preservatives   Additive = "4"
	Sweeteners      Additive = "5"
	Antioxidant     Additive = "7"
	FlavorEnhancers Additive = "8"
	Sulphurated     Additive = "9"
	Blackened       Additive = "10"
	Waxed           Additive = "11"
	Phosphate       Additive = "12"
	Phenylalanine   Additive = "13"
	Coating         Additive = "30"
)

var additiveOrder = fmap.Order(
	Color, Caffeine, Preservatives, Sweeteners, Antioxidant, FlavorEnhancers, Sulphurated, Blackened, Waxed, Phosphate, Phenylalanine,
	Coating,
)

func (a Additive) Cmp(other Additive) int {
	return additiveOrder[a] - additiveOrder[other]
}

func (a Additive) Known() bool {
	return additiveOrder.Has(a)
}

func (a Additive) Normalize() (Additive, bool) {
	key, _, ok := additiveOrder.Get(a)
	return key, ok
}

func (a Additive) ENString() string {
	return additivesEN[a]
}
func (a Additive) ENHTML() template.HTML {
	return template.HTML("<a class='annot' href='#add-" + string(a) + "' title='" + a.ENString() + "'>" + string(a) + "</a>")
}

func (a Additive) DEString() string {
	return additiveDE[a]
}

func (a Additive) DEHTML() template.HTML {
	return template.HTML("<a class='annot' href='#add-" + string(a) + "' title='" + a.DEString() + "'>" + string(a) + "</a>")
}

var additivesEN = map[Additive]string{
	Color:           "contains colour additives",
	Caffeine:        "contains caffeine",
	Preservatives:   "contains preservatives",
	Sweeteners:      "contains sweeteners",
	Antioxidant:     "contains antioxidant",
	FlavorEnhancers: "contains flavour enhancers",
	Sulphurated:     "sulphurated",
	Blackened:       "blackened",
	Waxed:           "waxed",
	Phosphate:       "contains phosphate",
	Phenylalanine:   "contains sweeteners = contains a source of phenylalanine",
	Coating:         "compound coating",
}

var additiveDE = map[Additive]string{
	Color:           "mit Farbstoff",
	Caffeine:        "mit Coffein",
	Preservatives:   "mit Konservierungsstoff",
	Sweeteners:      "mit Süßungsmittel",
	Antioxidant:     "mit Antioxidationsmittel",
	FlavorEnhancers: "mit Geschmacksverstärker",
	Sulphurated:     "geschwefelt",
	Blackened:       "geschwärzt",
	Waxed:           "gewachst",
	Phosphate:       "mit Phosphat",
	Phenylalanine:   "mit Süßungsmittel = enthält eine Phenylalaninquelle",
	Coating:         "mit Fettglasur",
}
