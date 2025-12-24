package annotations

import (
	"html/template"

	"github.com/tkw1536/faulunch/internal/fmap"
)

type Ingredient string

const (
	Vegetarian Ingredient = "V"

	Beef    Ingredient = "R"
	Poultry Ingredient = "G"
	Lamb    Ingredient = "L"
	FishI   Ingredient = "F"
	Pork    Ingredient = "S"
	Game    Ingredient = "W"

	Vegan      Ingredient = "veg"
	MensaVital Ingredient = "MV"
	Organic    Ingredient = "Bio"
	organicAlt Ingredient = "B" // alternate (newer) variant for Organic
	FishMSC    Ingredient = "MSC"

	Alcohol    Ingredient = "A"
	Glutenfree Ingredient = "Gf"
	CO2Neutral Ingredient = "CO2"
)

func (i *Ingredient) DoNormalize() {
	if *i == organicAlt {
		*i = Organic
	}
}

var ingredientOrder = fmap.Order(
	Vegetarian,
	Beef, Poultry, Lamb, FishI, Pork, Game,
	Vegan, MensaVital, Organic, FishMSC,
	Alcohol, Glutenfree, CO2Neutral,
)

var ingredientEN = map[Ingredient]string{
	Vegetarian: "vegetarian",

	Beef:    "beef",
	Poultry: "poultry",
	Lamb:    "lamb",
	FishI:   "fish",
	Pork:    "pork",
	Game:    "game",

	Vegan:      "vegan",
	MensaVital: "Cafeteria Vital",
	Organic:    "organic (certified by DE-ÖKO-006)",
	FishMSC:    "sustainable fish (certified by MSC - C - 51840)",

	Alcohol:    "with alcohol",
	Glutenfree: "gluten free",
	CO2Neutral: "CO2 Neutral",
}

var ingredientDE = map[Ingredient]string{
	Vegetarian: "Vegetarisch",

	Beef:    "Rind",
	Poultry: "Geflügel",
	Lamb:    "Lamm",
	FishI:   "Fisch",
	Pork:    "Schwein",
	Game:    "Wild",

	Vegan:      "Vegan",
	MensaVital: "Mensa Vital",
	Organic:    "aus biologischem Anbau DE-ÖKO-006",
	FishMSC:    "zertifizierte nachhaltige Fischerei - MSC - C - 51840",

	Alcohol:    "mit Alkohol",
	Glutenfree: "Glutenfrei",
	CO2Neutral: "CO2 Neutral",
}

func (i Ingredient) Cmp(other Ingredient) int {
	return ingredientOrder[i] - ingredientOrder[other]
}

func (i Ingredient) Known() bool {
	return ingredientOrder.Has(i)
}

func (i Ingredient) Normalize() (Ingredient, bool) {
	key, _, ok := ingredientOrder.Get(i)
	return key, ok
}

func (i Ingredient) ENString() string {
	return ingredientEN[i]
}
func (i Ingredient) ENHTML() template.HTML {
	return template.HTML("<a class='annot' href='#ing-" + string(i) + "' title='" + i.ENString() + "'>" + string(i) + "</a>")
}
func (i Ingredient) ENDef() template.HTML {
	return template.HTML("<a class='annot' href='#ing-" + string(i) + "' title='" + i.ENString() + "'>" + i.ENString() + "</a>")
}

func (i Ingredient) DEString() string {
	return ingredientDE[i]
}
func (i Ingredient) DEHTML() template.HTML {
	return template.HTML("<a class='annot' href='#ing-" + string(i) + "' title='" + i.DEString() + "'>" + string(i) + "</a>")
}
func (i Ingredient) DEDef() template.HTML {
	return template.HTML("<a class='annot' href='#ing-" + string(i) + "' title='" + i.DEString() + "'>" + i.DEString() + "</a>")
}
