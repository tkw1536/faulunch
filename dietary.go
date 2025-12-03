package faulunch

func (m *MenuItem) getDietaryCategory() DietaryCategory {
	ingredients := m.IngredientAnnotations.Data()

	// vegan
	for _, ing := range ingredients {
		if ing == Vegan {
			return DietaryCategoryVegan
		}
	}

	// vegetarian
	for _, ing := range ingredients {
		if ing == Vegetarian {
			return DietaryCategoryVegetarian
		}
	}

	// check if we have fish
	hasFish := false
	for _, ing := range ingredients {
		if ing == FishI {
			hasFish = true
			break
		}
	}
	if !hasFish {
		for _, allergen := range m.AllergenAnnotations.Data() {
			if allergen == Fish {
				hasFish = true
				break
			}
		}
	}
	if hasFish {
		return DietaryCategoryFish
	}

	// default to meat
	return DietaryCategoryMeat
}

// DietaryCategory represents the dietary category of a menu item.
type DietaryCategory string

// Different dietary categories
const (
	DietaryCategoryMeat       DietaryCategory = "meat"
	DietaryCategoryFish       DietaryCategory = "fish"
	DietaryCategoryVegetarian DietaryCategory = "vegetarian"
	DietaryCategoryVegan      DietaryCategory = "vegan"
)

func (d DietaryCategory) IsRestricted() bool {
	return d != "" && d != DietaryCategoryMeat
}

func (d DietaryCategory) ENString() string {
	switch d {
	case DietaryCategoryMeat:
		return "Meat"
	case DietaryCategoryFish:
		return "Fish"
	case DietaryCategoryVegetarian:
		return "Vegetarian"
	case DietaryCategoryVegan:
		return "Vegan"
	}
	panic("unknown dietary category")
}

func (d DietaryCategory) DEString() string {
	switch d {
	case DietaryCategoryMeat:
		return "Fleisch"
	case DietaryCategoryFish:
		return "Fisch"
	case DietaryCategoryVegetarian:
		return "Vegetarisch"
	case DietaryCategoryVegan:
		return "Vegan"
	}
	panic("unknown dietary category")
}
