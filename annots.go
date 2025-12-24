//spellchecker:words faulunch
package faulunch

//spellchecker:words html template regexp strings github zerolog faulunch internal fmap
import (
	"html/template"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal"
	"github.com/tkw1536/faulunch/internal/annotations"
)

var annotationPattern = regexp.MustCompile(`\([^)\s]+\)`)

func (item *MenuItem) extractAnnotations(logger *zerolog.Logger) {
	additives := make(map[annotations.Additive]struct{})
	allergens := make(map[annotations.Allergen]struct{})
	ingredients := make(map[annotations.Ingredient]struct{}, len(item.Piktogramme.Data()))

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

	internal.SetJSONData(&item.AdditiveAnnotations, internal.SortedKeysOf(additives, func(a, b annotations.Additive) int { return a.Cmp(b) }))
	internal.SetJSONData(&item.AllergenAnnotations, internal.SortedKeysOf(allergens, func(a, b annotations.Allergen) int { return a.Cmp(b) }))
	internal.SetJSONData(&item.IngredientAnnotations, internal.SortedKeysOf(ingredients, func(a, b annotations.Ingredient) int { return a.Cmp(b) }))
}

// RenderAnnotations renders annotations in the provided text
func (menu *MenuItem) renderAnnotations(logger *zerolog.Logger, text string, english bool, additives map[annotations.Additive]struct{}, allergens map[annotations.Allergen]struct{}, ingredients map[annotations.Ingredient]struct{}) template.HTML {
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
		if annotations.Additive(c).Known() || annotations.Allergen(c).Known() || annotations.Ingredient(c).Known() {
			return true
		}
	}
	return false
}

// renders and extracts a single annotation
func (menu *MenuItem) renderAnnot(logger *zerolog.Logger, annot string, english bool, additives map[annotations.Additive]struct{}, allergens map[annotations.Allergen]struct{}, ingredients map[annotations.Ingredient]struct{}) template.HTML {
	{
		add := annotations.Additive(annot)
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

		all := annotations.Allergen(annot)
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
		ing := annotations.Ingredient(annot)
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

var pictogramRegexp = regexp.MustCompile(regexp.QuoteMeta("https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/icons/") + `([^\.]+)` + regexp.QuoteMeta(".png"))

// ParseIngredients parses ingredients from a list of pictograms
func (menu *MenuItem) parseIngredients(s string, logger *zerolog.Logger) []annotations.Ingredient {
	ingredients := make(map[annotations.Ingredient]struct{})
	for _, match := range pictogramRegexp.FindAllStringSubmatch(s, -1) {
		ing := annotations.Ingredient(match[1])
		ing.DoNormalize()
		if !ing.Known() {
			logger.Error().Str("ingredient", match[1]).Time("day", menu.Day.Time()).Str("location", string(menu.Location)).Msg("Unknown Ingredient")
			continue
		}
		ingredients[ing] = struct{}{}
	}

	ings := internal.SortedKeysOf(ingredients, func(a, b annotations.Ingredient) int { return a.Cmp(b) })
	return ings
}
