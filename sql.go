//spellchecker:words faulunch
package faulunch

//spellchecker:words html template strconv strings github zerolog gorm datatypes
import (
	"html/template"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"gorm.io/datatypes"
)

// MenuItem represents a single item on a menu
type MenuItem struct {
	ID uint `gorm:"primaryKey" json:"-"`

	Day      Day      `gorm:"index" json:"-"` // the day this item is for
	Location Location `gorm:"index" json:"-"` // the location this item is for

	Category   string `gorm:"index"` // line this item is in
	CategoryEN string // english translation of category

	TitleDE string // title of this item in english
	TitleEN string // title of this item in german

	DescriptionDE string // description of this item (german)
	DescriptionEN string // description of this item (english)

	BeilagenDE string // sides (de)
	BeilagenEN string // sides (en)

	Preis1 LPrice // price (student)
	Preis2 LPrice // price (employee)
	Preis3 LPrice // price (guest)

	Piktogramme   datatypes.JSONType[[]Ingredient]
	Kj            LFloat
	Kcal          LFloat
	Fett          LFloat
	Gesfett       LFloat
	Kh            LFloat
	Zucker        LFloat
	Ballaststoffe LFloat
	Eiweiss       LFloat
	Salz          LFloat

	GlutenFree bool // is this gluten free?

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
	m.extractGlutenFree(logger)
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

func (m *MenuItem) extractGlutenFree(logger *zerolog.Logger) {
	m.GlutenFree = m.isGlutenFree()
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

type LPrice float64

func (lp LPrice) DEString() string {
	return strings.ReplaceAll(lp.ENString(), ".", ",")
}
func (lp LPrice) ENString() string {
	return strconv.FormatFloat(float64(lp), 'f', 2, 64)
}

// LFloat represents a localized float
type LFloat float64

func (lf LFloat) DEString() string {
	return strings.ReplaceAll(lf.ENString(), ".", ",")
}

func (lf LFloat) ENString() string {
	value := strconv.FormatFloat(float64(lf), 'f', 5, 64)
	return strings.TrimSuffix(strings.TrimRight(value, "0"), ".")
}
