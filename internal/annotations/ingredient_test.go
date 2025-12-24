package annotations_test

import (
	"testing"

	"github.com/tkw1536/faulunch/internal/annotations"
)

func TestIngredient_DoNormalize(t *testing.T) {
	tests := []struct {
		name  string
		input annotations.Ingredient
		want  annotations.Ingredient
	}{
		{name: "B becomes Bio", input: "B", want: annotations.Organic},
		{name: "Bio stays Bio", input: annotations.Organic, want: annotations.Organic},
		{name: "other unchanged", input: annotations.Vegetarian, want: annotations.Vegetarian},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input
			got.DoNormalize()
			if got != tt.want {
				t.Errorf("DoNormalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngredient_LocalizedString(t *testing.T) {
	tests := []struct {
		name   string
		i      annotations.Ingredient
		wantEN string
		wantDE string
	}{
		{name: "Vegetarian", i: annotations.Vegetarian, wantEN: "vegetarian", wantDE: "Vegetarisch"},
		{name: "Beef", i: annotations.Beef, wantEN: "beef", wantDE: "Rind"},
		{name: "Poultry", i: annotations.Poultry, wantEN: "poultry", wantDE: "Geflügel"},
		{name: "Lamb", i: annotations.Lamb, wantEN: "lamb", wantDE: "Lamm"},
		{name: "FishI", i: annotations.FishI, wantEN: "fish", wantDE: "Fisch"},
		{name: "Pork", i: annotations.Pork, wantEN: "pork", wantDE: "Schwein"},
		{name: "Game", i: annotations.Game, wantEN: "game", wantDE: "Wild"},
		{name: "Vegan", i: annotations.Vegan, wantEN: "vegan", wantDE: "Vegan"},
		{name: "MensaVital", i: annotations.MensaVital, wantEN: "Cafeteria Vital", wantDE: "Mensa Vital"},
		{name: "Organic", i: annotations.Organic, wantEN: "organic (certified by DE-ÖKO-006)", wantDE: "aus biologischem Anbau DE-ÖKO-006"},
		{name: "FishMSC", i: annotations.FishMSC, wantEN: "sustainable fish (certified by MSC - C - 51840)", wantDE: "zertifizierte nachhaltige Fischerei - MSC - C - 51840"},
		{name: "Alcohol", i: annotations.Alcohol, wantEN: "with alcohol", wantDE: "mit Alkohol"},
		{name: "Glutenfree", i: annotations.Glutenfree, wantEN: "gluten free", wantDE: "Glutenfrei"},
		{name: "CO2Neutral", i: annotations.CO2Neutral, wantEN: "CO2 Neutral", wantDE: "CO2 Neutral"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.ENString(); got != tt.wantEN {
				t.Errorf("ENString() = %v, want %v", got, tt.wantEN)
			}
			if got := tt.i.DEString(); got != tt.wantDE {
				t.Errorf("DEString() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}

func TestIngredient_LocalizedHTML(t *testing.T) {
	tests := []struct {
		name   string
		i      annotations.Ingredient
		wantEN string
		wantDE string
	}{
		{
			name:   "Vegetarian",
			i:      annotations.Vegetarian,
			wantEN: "<a class='annot' href='#ing-V' title='vegetarian'>V</a>",
			wantDE: "<a class='annot' href='#ing-V' title='Vegetarisch'>V</a>",
		},
		{
			name:   "Vegan",
			i:      annotations.Vegan,
			wantEN: "<a class='annot' href='#ing-veg' title='vegan'>veg</a>",
			wantDE: "<a class='annot' href='#ing-veg' title='Vegan'>veg</a>",
		},
		{
			name:   "Organic",
			i:      annotations.Organic,
			wantEN: "<a class='annot' href='#ing-Bio' title='organic (certified by DE-ÖKO-006)'>Bio</a>",
			wantDE: "<a class='annot' href='#ing-Bio' title='aus biologischem Anbau DE-ÖKO-006'>Bio</a>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.i.ENHTML()); got != tt.wantEN {
				t.Errorf("ENHTML() = %v, want %v", got, tt.wantEN)
			}
			if got := string(tt.i.DEHTML()); got != tt.wantDE {
				t.Errorf("DEHTML() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}

func TestIngredient_LocalizedDef(t *testing.T) {
	tests := []struct {
		name   string
		i      annotations.Ingredient
		wantEN string
		wantDE string
	}{
		{
			name:   "Vegetarian",
			i:      annotations.Vegetarian,
			wantEN: "<a class='annot' href='#ing-V' title='vegetarian'>vegetarian</a>",
			wantDE: "<a class='annot' href='#ing-V' title='Vegetarisch'>Vegetarisch</a>",
		},
		{
			name:   "Vegan",
			i:      annotations.Vegan,
			wantEN: "<a class='annot' href='#ing-veg' title='vegan'>vegan</a>",
			wantDE: "<a class='annot' href='#ing-veg' title='Vegan'>Vegan</a>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.i.ENDef()); got != tt.wantEN {
				t.Errorf("ENDef() = %v, want %v", got, tt.wantEN)
			}
			if got := string(tt.i.DEDef()); got != tt.wantDE {
				t.Errorf("DEDef() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}
