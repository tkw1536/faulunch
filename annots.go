package faulunch

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal"
	"github.com/tkw1536/faulunch/internal/fmap"
)

var annotationPattern = regexp.MustCompile(`\([^)\s]+\)`)

func (item *MenuItem) extractAnnotations(logger *zerolog.Logger) {
	additives := make(map[Additive]struct{})
	allergens := make(map[Allergen]struct{})
	ingredients := make(map[Ingredient]struct{}, len(item.Piktogramme.Data()))

	for _, ing := range item.Piktogramme.Data() {
		ingredients[ing] = struct{}{}
	}

	item.HTMLTitleDE = item.renderAnnotations(logger, item.TitleDE, false, additives, allergens, ingredients)
	item.HTMLTitleEN = item.renderAnnotations(logger, item.TitleEN, true, additives, allergens, ingredients)

	item.HTMLDescriptionDE = item.renderAnnotations(logger, item.DescriptionDE, false, additives, allergens, ingredients)
	item.HTMLDescriptionEN = item.renderAnnotations(logger, item.DescriptionEN, true, additives, allergens, ingredients)

	item.HTMLBeilagenDE = item.renderAnnotations(logger, item.BeilagenDE, false, additives, allergens, ingredients)
	item.HTMLBeilagenEN = item.renderAnnotations(logger, item.BeilagenEN, true, additives, allergens, ingredients)

	// store all the additive and ingredient data
	// then sort it for convenience

	internal.SetJSONData(&item.AdditiveAnnotations, internal.SortedKeysOf(additives, func(a, b Additive) int { return a.Cmp(b) }))
	internal.SetJSONData(&item.AllergenAnnotations, internal.SortedKeysOf(allergens, func(a, b Allergen) int { return a.Cmp(b) }))
	internal.SetJSONData(&item.IngredientAnnotations, internal.SortedKeysOf(ingredients, func(a, b Ingredient) int { return a.Cmp(b) }))
}

// RenderAnnotations renders annotations in the provided text
func (menu *MenuItem) renderAnnotations(logger *zerolog.Logger, text string, english bool, additives map[Additive]struct{}, allergens map[Allergen]struct{}, ingredients map[Ingredient]struct{}) template.HTML {
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
		annots := strings.FieldsFunc(matches[i], func(r rune) bool { return r == ',' || r == '.' })
		annots = fixAnnotTypos(annots)
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
		for _, annot := range annots {
			buffer = append(buffer, string(menu.renderAnnot(logger, annot, english, additives, allergens, ingredients)))
		}

		builder.WriteString(strings.Join(buffer, ", "))
		builder.WriteString("</span>")
	}

	return template.HTML(builder.String())
}

// fixes typos in the annotations
func fixAnnotTypos(annots []string) []string {
	fix := make([]string, 0, len(annots))
	for _, a := range annots {
		switch a {
		case "Vegan":
			fix = append(fix, "veg")
		case "EiEi", "Egg":
			fix = append(fix, "Ei")
		case "Mi7":
			fix = append(fix, "Mi", "7")
		case "Sel1":
			fix = append(fix, "Sel", "1")
		case "RWz":
			fix = append(fix, "R", "Wz")
		case "Sul": // Not sure about this one
			fix = append(fix, "Su")
		case "VWz":
			fix = append(fix, "V", "Wz")
		case "SelGe":
			fix = append(fix, "Sel", "Ge")
		case "SuGe":
			fix = append(fix, "Su", "Ge")
		case "Wzel":
			fix = append(fix, "Wz")
		case "Sun":
			fix = append(fix, "So")
		case "Ma":
			fix = append(fix, "Man")
		case "Wed":
			fix = append(fix, "Mi")
		case "Se":
			fix = append(fix, "Wa")
		case "3", "cond":
			// these are mismatches and no longer used
		default:
			fix = append(fix, a)
		}
	}
	return fix
}

// validMatches checks if at least one annotation inside the match is valid
func anyValidAnnot(matches ...string) bool {
	for _, c := range matches {
		if Additive(c).Known() || Allergen(c).Known() || Ingredient(c).Known() {
			return true
		}
	}
	return false
}

// renders and extracts a single annotation
func (menu *MenuItem) renderAnnot(logger *zerolog.Logger, annot string, english bool, additives map[Additive]struct{}, allergens map[Allergen]struct{}, ingredients map[Ingredient]struct{}) template.HTML {
	{
		add := Additive(annot)
		if add, ok := add.Normalize(); ok {
			additives[add] = struct{}{}
			if english {
				return add.ENHTML()
			} else {
				return add.DEHTML()
			}
		}
	}

	{

		all := Allergen(annot)
		if all, ok := all.Normalize(); ok {
			allergens[all] = struct{}{}
			if english {
				return all.ENHTML()
			} else {
				return all.DEHTML()
			}
		}
	}

	{
		ing := Ingredient(annot)
		if ing, ok := ing.Normalize(); ok {
			ingredients[ing] = struct{}{}
			if english {
				return ing.ENHTML()
			} else {
				return ing.DEHTML()
			}
		}
	}

	logger.Error().Str("annot", annot).Int("day", int(menu.Day)).Str("location", string(menu.Location)).Bool("english", english).Msg("Unknown annotation")

	return template.HTML(template.HTMLEscapeString(annot))
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
	Waxed           Additive = "11"
	Phosphate       Additive = "12"
	Phenylalanine   Additive = "13"
	Coating         Additive = "30"
)

var additiveOrder = order(
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

var allergenOrder = order(
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

type Ingredient string

var pictogramRegexp = regexp.MustCompile(regexp.QuoteMeta("https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/icons/") + `([^\.]+)` + regexp.QuoteMeta(".png"))

// ParseIngredients parses ingredients from a list of pictograms
func (menu *MenuItem) parseIngredients(s string, logger *zerolog.Logger) []Ingredient {
	ingredients := make(map[Ingredient]struct{})
	for _, match := range pictogramRegexp.FindAllStringSubmatch(s, -1) {
		ing := Ingredient(match[1])
		if !ing.Known() {
			logger.Error().Str("ingredient", match[1]).Time("day", menu.Day.Time()).Str("location", string(menu.Location)).Msg("Unknown Ingredient")
			continue
		}
		ingredients[ing] = struct{}{}
	}

	ings := internal.SortedKeysOf(ingredients, func(a, b Ingredient) int { return a.Cmp(b) })
	return ings
}

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
	FishMSC    Ingredient = "MSC"

	Alcohol    Ingredient = "A"
	Glutenfree Ingredient = "Gf"
	CO2Neutral Ingredient = "CO2"
)

var ingredientOrder = order(
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

func order[T ~string](values ...T) fmap.FMap[T, int] {
	m := make(fmap.FMap[T, int], len(values))
	for index, item := range values {
		m.Add(item, index)
	}
	return m
}
