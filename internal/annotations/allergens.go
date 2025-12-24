package annotations

import (
	"html/template"

	"github.com/tkw1536/faulunch/internal/fmap"
)

type Allergen string

func (a Allergen) Cmp(other Allergen) int {
	return allergenOrder[a] - allergenOrder[other]
}

func (a Allergen) Known() bool {
	return allergenOrder.Has(a)
}

func (a Allergen) Normalize() (Allergen, bool) {
	key, _, ok := allergenOrder.Get(a)
	return key, ok
}

func (a Allergen) ENString() string {
	return allergensEN[a]
}
func (a Allergen) ENHTML() template.HTML {
	return template.HTML("<a class='annot' href='#all-" + string(a) + "' title='" + a.ENString() + "'>" + string(a) + "</a>")
}

func (a Allergen) DEString() string {
	return allergensDE[a]
}
func (a Allergen) DEHTML() template.HTML {
	return template.HTML("<a class='annot' href='#all-" + string(a) + "' title='" + a.DEString() + "'>" + string(a) + "</a>")
}

const (
	Wheat         Allergen = "Wz"
	Rye           Allergen = "Ro"
	Barley        Allergen = "Ge"
	Oats          Allergen = "Hf"
	Crustaceans   Allergen = "Kr"
	Eggs          Allergen = "Ei"
	Fish          Allergen = "Fi"
	Peanuts       Allergen = "Er"
	Soybeans      Allergen = "So"
	Milk          Allergen = "Mi"
	Almonds       Allergen = "Man"
	HazelNuts     Allergen = "Hs"
	WalNuts       Allergen = "Wa"
	CashewNuts    Allergen = "Ka"
	PecanNuts     Allergen = "Pe"
	BrazilNuts    Allergen = "Pa"
	Pistachios    Allergen = "Pi"
	MacadamiaNuts Allergen = "Mac"
	Celeriac      Allergen = "Sel"
	Mustard       Allergen = "Sen"
	Sesame        Allergen = "Ses"
	Sulphur       Allergen = "Su"
	Lupines       Allergen = "Lu"
	Mollusca      Allergen = "We"
)

var allergenOrder = fmap.Order(
	Wheat, Rye, Barley, Oats, Crustaceans, Eggs, Fish, Peanuts, Soybeans, Milk,
	Almonds, HazelNuts, WalNuts, CashewNuts, PecanNuts, BrazilNuts, Pistachios, MacadamiaNuts, Celeriac, Mustard,
	Sesame, Sulphur, Lupines, Mollusca,
)

var allergensEN = map[Allergen]string{
	Wheat:         "cereals containing gluten wheat (spelt, kamut)",
	Rye:           "cereals containing gluten rye",
	Barley:        "cereals containing gluten barley",
	Oats:          "cereals containing gluten oats",
	Crustaceans:   "contains crustaceans",
	Eggs:          "eggs",
	Fish:          "fish",
	Peanuts:       "peanuts",
	Soybeans:      "soybeans",
	Milk:          "milk/lactose",
	Almonds:       "almonds",
	HazelNuts:     "hazelnuts",
	WalNuts:       "walnuts",
	CashewNuts:    "cashew nuts",
	PecanNuts:     "pecan nuts",
	BrazilNuts:    "brazil nuts",
	Pistachios:    "pistachios",
	MacadamiaNuts: "macadamia nuts",
	Celeriac:      "celeriac",
	Mustard:       "mustard",
	Sesame:        "sesame",
	Sulphur:       "sulphur dioxide and sulphites",
	Lupines:       "lupines",
	Mollusca:      "mollusca",
}

var allergensDE = map[Allergen]string{
	Wheat:         "glutenhaltiges Getreide Weizen (Dinkel, Kamut)",
	Rye:           "glutenhaltiges Getreide Roggen",
	Barley:        "glutenhaltiges Getreide Gerste",
	Oats:          "glutenhaltiges Getreide Hafer",
	Crustaceans:   "Krebstiere",
	Eggs:          "Eier",
	Fish:          "Fisch",
	Peanuts:       "Erdnüsse",
	Soybeans:      "Sojabohnen",
	Milk:          "Milch/Laktose",
	Almonds:       "Schalenfrüchte Mandeln",
	HazelNuts:     "Schalenfrüchte Haselnüsse",
	WalNuts:       "Schalenfrüchte Walnüsse",
	CashewNuts:    "Schalenfrüchte Kaschu(Cashew)nüsse",
	PecanNuts:     "Schalenfrüchte Pekannüsse",
	BrazilNuts:    "Schalenfrüchte Paranüsse",
	Pistachios:    "Schalenfrüchte Pistazien",
	MacadamiaNuts: "Schalenfrüchte Macadamianüsse",
	Celeriac:      "Sellerie",
	Mustard:       "Senf",
	Sesame:        "Sesam",
	Sulphur:       "Schwefeldioxid und Sulfite",
	Lupines:       "Lupinen",
	Mollusca:      "Weichtiere",
}
