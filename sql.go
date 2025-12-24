//spellchecker:words faulunch
package faulunch

//spellchecker:words html template strconv strings github zerolog gorm datatypes
import (
	"html/template"
	"strings"

	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal/ltime"
	"github.com/tkw1536/faulunch/internal/types"
	"gorm.io/datatypes"
)

// MenuItem represents a single item on a menu
type MenuItem struct {
	ID uint `gorm:"primaryKey" json:"-"`

	Day      ltime.Day `gorm:"index" json:"-"` // the day this item is for
	Location Location  `gorm:"index" json:"-"` // the location this item is for

	Category   string `gorm:"index"` // line this item is in
	CategoryEN string // english translation of category

	TitleDE string // title of this item in english
	TitleEN string // title of this item in german

	DescriptionDE string // description of this item (german)
	DescriptionEN string // description of this item (english)

	BeilagenDE string // sides (de)
	BeilagenEN string // sides (en)

	Preis1 types.LPrice // price (student)
	Preis2 types.LPrice // price (employee)
	Preis3 types.LPrice // price (guest)

	Piktogramme   datatypes.JSONType[[]Ingredient]
	Kj            types.LFloat
	Kcal          types.LFloat
	Fett          types.LFloat
	Gesfett       types.LFloat
	Kh            types.LFloat
	Zucker        types.LFloat
	Ballaststoffe types.LFloat
	Eiweiss       types.LFloat
	Salz          types.LFloat

	GlutenFree      bool            // is this gluten free?
	DietaryCategory DietaryCategory // the dietary category of this item

	// Annotations properly replaced with <span class='#type'> and inside <sup>s
	HTMLTitleDE       template.HTML
	HTMLTitleEN       template.HTML
	HTMLDescriptionDE template.HTML
	HTMLDescriptionEN template.HTML
	HTMLBeilagenDE    template.HTML
	HTMLBeilagenEN    template.HTML

	AllergenAnnotations   datatypes.JSONType[[]Allergen]
	AdditiveAnnotations   datatypes.JSONType[[]Additive]
	IngredientAnnotations datatypes.JSONType[[]Ingredient]
}

var categoryTranslations = map[string]string{
	"Essen":          "Meal",
	"Aktionsessen":   "Special Meal",
	"Aktion":         "Special",
	"Suppe":          "Soup",
	"Suppen":         "Soups",
	"SB-Theke":       "Self-Service Counter",
	"Tagesangebot":   "Daily Special",
	"Tipp des Tages": "Tip Of The Day",
}

func (m *MenuItem) UpdateComputedFields(logger *zerolog.Logger) {
	m.translateCategoryNames(logger)
	m.extractAnnotations(logger)
	m.extractGlutenFree()
	m.extractDietaryCategory()
}

func (m *MenuItem) translateCategoryNames(logger *zerolog.Logger) {
	if complete, ok := categoryTranslations[m.Category]; ok {
		m.CategoryEN = complete
		return
	}

	fields := strings.Fields(m.Category)
	for i, field := range fields {
		trans, ok := categoryTranslations[field]
		if ok {
			fields[i] = trans
		} else if !isOnlyDigits(field) {
			logger.Debug().Str("part", field).Msg("untranslatable category part")
		}
	}

	m.CategoryEN = strings.Join(fields, " ")
}

func (m *MenuItem) extractGlutenFree() {
	m.GlutenFree = m.isGlutenFree()
}

func (m *MenuItem) extractDietaryCategory() {
	m.DietaryCategory = m.getDietaryCategory()
}

func (m MenuItem) isGlutenFree() bool {
	for _, allergen := range m.AllergenAnnotations.Data() {
		if allergen == Wheat || allergen == Rye || allergen == Barley || allergen == Oats {
			return false
		}
	}
	return true
}

func isOnlyDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (m MenuItem) Ingredients() []Ingredient {
	return m.Piktogramme.Data()
}

func (m MenuItem) catScore() int {
	switch {
	case strings.HasPrefix(m.Category, "Essen "):
		return -3
	case strings.HasPrefix(m.Category, "Aktionsessen "):
		return -2
	case strings.HasPrefix(m.Category, "Suppe "):
		return -1
	default:
		return 0
	}
}
func (m MenuItem) Cmp(other MenuItem) int {
	ours, theirs := m.catScore(), other.catScore()
	if ours < theirs {
		return -1
	} else if ours > theirs {
		return 1
	}
	return strings.Compare(m.Category, other.Category)
}
