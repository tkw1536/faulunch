package annotations_test

import (
	"testing"

	"github.com/tkw1536/faulunch/internal/annotations"
)

func TestAdditive_LocalizedString(t *testing.T) {
	tests := []struct {
		name   string
		a      annotations.Additive
		wantEN string
		wantDE string
	}{
		{name: "Color", a: annotations.Color, wantEN: "contains colour additives", wantDE: "mit Farbstoff"},
		{name: "Caffeine", a: annotations.Caffeine, wantEN: "contains caffeine", wantDE: "mit Coffein"},
		{name: "Preservatives", a: annotations.Preservatives, wantEN: "contains preservatives", wantDE: "mit Konservierungsstoff"},
		{name: "Sweeteners", a: annotations.Sweeteners, wantEN: "contains sweeteners", wantDE: "mit Süßungsmittel"},
		{name: "Antioxidant", a: annotations.Antioxidant, wantEN: "contains antioxidant", wantDE: "mit Antioxidationsmittel"},
		{name: "FlavorEnhancers", a: annotations.FlavorEnhancers, wantEN: "contains flavour enhancers", wantDE: "mit Geschmacksverstärker"},
		{name: "Sulphurated", a: annotations.Sulphurated, wantEN: "sulphurated", wantDE: "geschwefelt"},
		{name: "Blackened", a: annotations.Blackened, wantEN: "blackened", wantDE: "geschwärzt"},
		{name: "Waxed", a: annotations.Waxed, wantEN: "waxed", wantDE: "gewachst"},
		{name: "Phosphate", a: annotations.Phosphate, wantEN: "contains phosphate", wantDE: "mit Phosphat"},
		{name: "Phenylalanine", a: annotations.Phenylalanine, wantEN: "contains sweeteners = contains a source of phenylalanine", wantDE: "mit Süßungsmittel = enthält eine Phenylalaninquelle"},
		{name: "Coating", a: annotations.Coating, wantEN: "compound coating", wantDE: "mit Fettglasur"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.ENString(); got != tt.wantEN {
				t.Errorf("ENString() = %v, want %v", got, tt.wantEN)
			}
			if got := tt.a.DEString(); got != tt.wantDE {
				t.Errorf("DEString() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}

func TestAdditive_LocalizedHTML(t *testing.T) {
	tests := []struct {
		name   string
		a      annotations.Additive
		wantEN string
		wantDE string
	}{
		{
			name:   "Color",
			a:      annotations.Color,
			wantEN: "<a class='annot' href='#add-1' title='contains colour additives'>1</a>",
			wantDE: "<a class='annot' href='#add-1' title='mit Farbstoff'>1</a>",
		},
		{
			name:   "Coating",
			a:      annotations.Coating,
			wantEN: "<a class='annot' href='#add-30' title='compound coating'>30</a>",
			wantDE: "<a class='annot' href='#add-30' title='mit Fettglasur'>30</a>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.a.ENHTML()); got != tt.wantEN {
				t.Errorf("ENHTML() = %v, want %v", got, tt.wantEN)
			}
			if got := string(tt.a.DEHTML()); got != tt.wantDE {
				t.Errorf("DEHTML() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}
