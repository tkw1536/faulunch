package annotations_test

import (
	"testing"

	"github.com/tkw1536/faulunch/internal/annotations"
)

func TestAllergen_LocalizedString(t *testing.T) {
	tests := []struct {
		name   string
		a      annotations.Allergen
		wantEN string
		wantDE string
	}{
		{name: "Wheat", a: annotations.Wheat, wantEN: "cereals containing gluten wheat (spelt, kamut)", wantDE: "glutenhaltiges Getreide Weizen (Dinkel, Kamut)"},
		{name: "Rye", a: annotations.Rye, wantEN: "cereals containing gluten rye", wantDE: "glutenhaltiges Getreide Roggen"},
		{name: "Barley", a: annotations.Barley, wantEN: "cereals containing gluten barley", wantDE: "glutenhaltiges Getreide Gerste"},
		{name: "Oats", a: annotations.Oats, wantEN: "cereals containing gluten oats", wantDE: "glutenhaltiges Getreide Hafer"},
		{name: "Crustaceans", a: annotations.Crustaceans, wantEN: "contains crustaceans", wantDE: "Krebstiere"},
		{name: "Eggs", a: annotations.Eggs, wantEN: "eggs", wantDE: "Eier"},
		{name: "Fish", a: annotations.Fish, wantEN: "fish", wantDE: "Fisch"},
		{name: "Peanuts", a: annotations.Peanuts, wantEN: "peanuts", wantDE: "Erdnüsse"},
		{name: "Soybeans", a: annotations.Soybeans, wantEN: "soybeans", wantDE: "Sojabohnen"},
		{name: "Milk", a: annotations.Milk, wantEN: "milk/lactose", wantDE: "Milch/Laktose"},
		{name: "Almonds", a: annotations.Almonds, wantEN: "almonds", wantDE: "Schalenfrüchte Mandeln"},
		{name: "HazelNuts", a: annotations.HazelNuts, wantEN: "hazelnuts", wantDE: "Schalenfrüchte Haselnüsse"},
		{name: "WalNuts", a: annotations.WalNuts, wantEN: "walnuts", wantDE: "Schalenfrüchte Walnüsse"},
		{name: "CashewNuts", a: annotations.CashewNuts, wantEN: "cashew nuts", wantDE: "Schalenfrüchte Kaschu(Cashew)nüsse"},
		{name: "PecanNuts", a: annotations.PecanNuts, wantEN: "pecan nuts", wantDE: "Schalenfrüchte Pekannüsse"},
		{name: "BrazilNuts", a: annotations.BrazilNuts, wantEN: "brazil nuts", wantDE: "Schalenfrüchte Paranüsse"},
		{name: "Pistachios", a: annotations.Pistachios, wantEN: "pistachios", wantDE: "Schalenfrüchte Pistazien"},
		{name: "MacadamiaNuts", a: annotations.MacadamiaNuts, wantEN: "macadamia nuts", wantDE: "Schalenfrüchte Macadamianüsse"},
		{name: "Celeriac", a: annotations.Celeriac, wantEN: "celeriac", wantDE: "Sellerie"},
		{name: "Mustard", a: annotations.Mustard, wantEN: "mustard", wantDE: "Senf"},
		{name: "Sesame", a: annotations.Sesame, wantEN: "sesame", wantDE: "Sesam"},
		{name: "Sulphur", a: annotations.Sulphur, wantEN: "sulphur dioxide and sulphites", wantDE: "Schwefeldioxid und Sulfite"},
		{name: "Lupines", a: annotations.Lupines, wantEN: "lupines", wantDE: "Lupinen"},
		{name: "Mollusca", a: annotations.Mollusca, wantEN: "mollusca", wantDE: "Weichtiere"},
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

func TestAllergen_LocalizedHTML(t *testing.T) {
	tests := []struct {
		name   string
		a      annotations.Allergen
		wantEN string
		wantDE string
	}{
		{
			name:   "Wheat",
			a:      annotations.Wheat,
			wantEN: "<a class='annot' href='#all-Wz' title='cereals containing gluten wheat (spelt, kamut)'>Wz</a>",
			wantDE: "<a class='annot' href='#all-Wz' title='glutenhaltiges Getreide Weizen (Dinkel, Kamut)'>Wz</a>",
		},
		{
			name:   "Milk",
			a:      annotations.Milk,
			wantEN: "<a class='annot' href='#all-Mi' title='milk/lactose'>Mi</a>",
			wantDE: "<a class='annot' href='#all-Mi' title='Milch/Laktose'>Mi</a>",
		},
		{
			name:   "MacadamiaNuts",
			a:      annotations.MacadamiaNuts,
			wantEN: "<a class='annot' href='#all-Mac' title='macadamia nuts'>Mac</a>",
			wantDE: "<a class='annot' href='#all-Mac' title='Schalenfrüchte Macadamianüsse'>Mac</a>",
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
