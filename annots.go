package faulunch

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var annotationPattern = regexp.MustCompile(`\([^)\s]+\)`)

// RenderAnnotations renders annotations in the provided text
func RenderAnnotations(text string, english bool) template.HTML {
	values := annotationPattern.Split(text, -1)
	for i, v := range values {
		values[i] = template.HTMLEscapeString(v)
	}

	matches := annotationPattern.FindAllString(text, -1)

	var builder strings.Builder
	var buffer []string
	for i, value := range values {

		// write the current non-match
		// and check if there is a following value
		builder.WriteString(value)
		if len(matches) <= i {
			continue
		}

		// trim off the leading and trailing bracket
		matches[i] = matches[i][1 : len(matches[i])-1]

		// find all the individual annotations
		// and check if there is at least one valid one
		annots := strings.Split(matches[i], ",")
		if !anyValidAnnot(annots...) {
			// no valid annotation => skip
			builder.WriteRune('(')
			builder.WriteString(matches[i])
			builder.WriteRune(')')
			continue
		}

		// replace all the valid annotations

		builder.WriteString("<span class='annot'>")

		buffer = buffer[:0]
		for _, annot := range strings.Split(matches[i], ",") {
			buffer = append(buffer, string(renderAnnot(annot, english)))
		}

		builder.WriteString(strings.Join(buffer, ", "))
		builder.WriteString("</span>")
	}

	return template.HTML(builder.String())
}

// validMatches checks if at least one annotation inside the match is valid
func anyValidAnnot(matches ...string) bool {
	for _, c := range matches {
		if validAnnot(c) {
			return true
		}
	}
	return false
}

// validAnnot checks if annot is valid
func validAnnot(annot string) bool {
	return Additive(annot).Known() || Allergen(annot).Known() || Ingredient(annot).Known()
}

func renderAnnot(annot string, english bool) template.HTML {
	{
		add := Additive(annot)
		if add.Known() {
			if english {
				return add.ENHTML()
			} else {
				return add.DEHTML()
			}
		}
	}

	{

		all := Allergen(annot)
		if all.Known() {
			if english {
				return all.ENHTML()
			} else {
				return all.DEHTML()
			}
		}
	}

	{
		ing := Ingredient(annot)
		if ing.Known() {
			if english {
				return ing.ENHTML()
			} else {
				return ing.DEHTML()
			}
		}
	}

	return template.HTML(template.HTMLEscapeString(annot))
}

// MenuAnnotations returns the annotations for the given menu items
func MenuAnnotations(items []MenuItem, logger *zerolog.Logger) ([]Additive, []Allergen, []Ingredient) {
	annots := make(map[string]struct{})
	for _, item := range items {
		extractAnnots(item.TitleDE, annots)
		extractAnnots(item.TitleEN, annots)
		extractAnnots(item.DescriptionDE, annots)
		extractAnnots(item.DescriptionEN, annots)
		extractAnnots(item.BeilagenDE, annots)
		extractAnnots(item.BeilagenEN, annots)

		for _, ing := range item.Ingredients() {
			annots[string(ing)] = struct{}{}
		}
	}

	return annotations(annots, logger)
}

func extractAnnots(source string, annots map[string]struct{}) {
	for _, value := range annotationPattern.FindAllString(source, -1) {
		value = value[1 : len(value)-1]
		values := strings.Split(value, ",")
		if !anyValidAnnot(values...) {
			continue
		}
		for _, annot := range values {
			annots[annot] = struct{}{}
		}
	}
}

func annotations(annots map[string]struct{}, logger *zerolog.Logger) (adds []Additive, alls []Allergen, ings []Ingredient) {
	// check if we have an additive or an allergen
	for annot := range annots {
		add := Additive(annot)
		if add.Known() {
			adds = append(adds, add)
			continue
		}

		all := Allergen(annot)
		if all.Known() {
			alls = append(alls, all)
			continue
		}

		ing := Ingredient(annot)
		if ing.Known() {
			ings = append(ings, ing)
			continue
		}

		logger.Error().Str("annotation", annot).Msg("Unknown annotation")
	}

	// sort the results
	slices.SortFunc(adds, func(a, b Additive) bool { return a.Less(b) })
	slices.SortFunc(alls, func(a, b Allergen) bool { return a.Less(b) })
	slices.SortFunc(ings, func(a, b Ingredient) bool { return a.Less(b) })

	// and done!
	return
}

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
	Phosphate       Additive = "12"
	Phenylalanine   Additive = "13"
	Coating         Additive = "30"
)

var additiveOrder = order(
	Color, Caffeine, Preservatives, Sweeteners, Antioxidant, FlavorEnhancers, Sulphurated, Blackened, Phosphate, Phenylalanine,
	Coating,
)

func (a Additive) Less(other Additive) bool {
	return additiveOrder[a] < additiveOrder[other]
}

func (a Additive) Known() bool {
	_, ok := additiveOrder[a]
	return ok
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
	Phosphate:       "contains phosphate",
	Phenylalanine:   "contains sweeteners, contains a source of phenylalanine",
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
	Phosphate:       "mit Phosphat",
	Phenylalanine:   "enthält eine Phenylalaninquelle",
	Coating:         "mit Fettglasur",
}

type Allergen string

func (a Allergen) Less(other Allergen) bool {
	return allergenOrder[a] < allergenOrder[other]
}

func (a Allergen) Known() bool {
	_, ok := allergenOrder[a]
	return ok
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

var allergenOrder = order(
	Wheat, Rye, Barley, Oats, Crustaceans, Eggs, Fish, Peanuts, Soybeans, Milk,
	Almonds, HazelNuts, WalNuts, CashewNuts, PecanNuts, BrazilNuts, Pistachios, MacadamiaNuts, Celeriac, Mustard,
	Sesame, Sulphur, Lupines, Mollusca,
)

var allergensEN = map[Allergen]string{
	Wheat:         "Cereals containing gluten wheat (spelt, kamut)",
	Rye:           "Cereals containing gluten rye",
	Barley:        "Cereals containing gluten barley",
	Oats:          "Cereals containing gluten oats",
	Crustaceans:   "contains Crustaceans",
	Eggs:          "Eggs",
	Fish:          "Fish",
	Peanuts:       "Peanuts",
	Soybeans:      "Soybeans",
	Milk:          "Milk/milk sugar",
	Almonds:       "Almonds",
	HazelNuts:     "Hazelnuts",
	WalNuts:       "Walnuts",
	CashewNuts:    "Cashew nuts",
	PecanNuts:     "Pecan nuts",
	BrazilNuts:    "Brazil nuts",
	Pistachios:    "Pistachios",
	MacadamiaNuts: "Macadamia nuts",
	Celeriac:      "Celeriac",
	Mustard:       "Mustard",
	Sesame:        "Sesame",
	Sulphur:       "Sulphur dioxide and sulphites",
	Lupines:       "Lupines",
	Mollusca:      "Mollusca",
}

var allergensDE = map[Allergen]string{
	Wheat:         "Weizen(Dinkel,Kamut)",
	Rye:           "Roggen",
	Barley:        "Gerste",
	Oats:          "Hafer",
	Crustaceans:   "Krebstiere",
	Eggs:          "Eier",
	Fish:          "Fisch",
	Peanuts:       "Erdnüsse",
	Soybeans:      "Sojabohnen",
	Milk:          "Milch/Laktose",
	Almonds:       "Mandeln",
	HazelNuts:     "Haselnüsse",
	WalNuts:       "Walnüsse",
	CashewNuts:    "Kaschu(Cashew)nüsse",
	PecanNuts:     "Pekannüsse",
	BrazilNuts:    "Paranüsse",
	Pistachios:    "Schalenfrüchte Pistazien",
	MacadamiaNuts: "Macadamianüsse",
	Celeriac:      "Sellerie",
	Mustard:       "Senf",
	Sesame:        "Sesam",
	Sulphur:       "Schwefeldioxid u. Sulfite",
	Lupines:       "Lupinen",
	Mollusca:      "Weichtiere",
}

type Ingredient string

var pictogramRegexp = regexp.MustCompile(regexp.QuoteMeta("https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/icons/") + `([^\.]+)` + regexp.QuoteMeta(".png"))

// ParseIngredients parses ingredients from a list of pictograms
func ParseIngredients(s string, logger *zerolog.Logger) []Ingredient {
	ingredients := make(map[Ingredient]struct{})
	for _, match := range pictogramRegexp.FindAllStringSubmatch(s, -1) {
		ing := Ingredient(match[1])
		if !ing.Known() {
			logger.Error().Str("ingredient", match[1]).Msg("Unknown Ingredient")
			continue
		}
		ingredients[ing] = struct{}{}
	}

	ings := maps.Keys(ingredients)
	slices.SortFunc(ings, func(a, b Ingredient) bool { return a.Less(b) })
	return ings
}

const (
	Pork       Ingredient = "S"
	Beef       Ingredient = "R"
	Poultry    Ingredient = "G"
	Lamb       Ingredient = "L"
	Game       Ingredient = "W"
	FishI      Ingredient = "F"
	Vegetarian Ingredient = "V"
	Vegan      Ingredient = "veg"
	Organic    Ingredient = "Bio"
	FishMSC    Ingredient = "MSC"
	Alcohol    Ingredient = "O" //  TODO: Is this the right annotation

	Glutenfree Ingredient = "Gf"
	MensaVital Ingredient = "MV"
	CO2Neutral Ingredient = "CO2"
)

var ingredientOrder = order(
	Pork, Beef, Poultry, Lamb, Game, FishI, Vegetarian, Vegan, Organic, FishMSC,
	Alcohol, MensaVital, CO2Neutral, Glutenfree,
)

var ingredientEN = map[Ingredient]string{
	Pork:       "Pork",
	Beef:       "Beef",
	Poultry:    "Poultry",
	Lamb:       "Lamb",
	Game:       "Game",
	FishI:      "Fish",
	Vegetarian: "Vegetarian",
	Vegan:      "Vegan",
	Organic:    "organic (certified by DE-ÖKO-006)",
	FishMSC:    "sustainable fish (certified by MSC - C - 51840)",
	Alcohol:    "with alcohol",

	Glutenfree: "Gluten Free",
	MensaVital: "Mensa Vital",
	CO2Neutral: "CO2 Neutral",
}

var ingredientDE = map[Ingredient]string{
	Pork:       "Schwein",
	Beef:       "Rind",
	Poultry:    "Geflügel",
	Lamb:       "Lamm",
	Game:       "Wild",
	FishI:      "Fisch",
	Vegetarian: "Vegetarisch",
	Vegan:      "Vegan",
	Organic:    "aus biologischem Anbau DE-ÖKO-006",
	FishMSC:    "zertifizierte nachhaltige Fischerei - MSC - C - 51840",
	Alcohol:    "mit Alkohol",

	Glutenfree: "Glutenfrei",
	MensaVital: "Mensa Vital",
	CO2Neutral: "CO2 Neutral",
}

func (i Ingredient) Less(other Ingredient) bool {
	return ingredientOrder[i] < ingredientOrder[other]
}

func (i Ingredient) Known() bool {
	_, ok := ingredientOrder[i]
	return ok
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

func order[T comparable](values ...T) map[T]int {
	m := make(map[T]int, len(values))
	for index, item := range values {
		m[item] = index
	}
	return m
}
