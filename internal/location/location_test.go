//spellchecker:words faulunch
package location_test

//spellchecker:words encoding json strings testing github faulunch internal location
import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tkw1536/faulunch/internal/location"
)

func TestLocation_Valid(t *testing.T) {
	tests := []struct {
		name     string
		location location.Location
		want     bool
	}{
		{name: "valid MensaSued", location: location.MensaSued, want: true},
		{name: "valid MensaInselschuett", location: location.MensaInselschuett, want: true},
		{name: "valid CafeteriaComeIn", location: location.CafeteriaComeIn, want: true},
		{name: "valid WohnanlageStPeter", location: location.WohnanlageStPeter, want: true},
		{name: "invalid empty", location: location.Location(""), want: false},
		{name: "invalid random", location: location.Location("random-location"), want: false},
		{name: "invalid typo", location: location.Location("mensa-suedd"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.location.Valid(); got != tt.want {
				t.Errorf("Location.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocations(t *testing.T) {
	locs := location.Locations()

	// Check that we have at least some locations
	if len(locs) == 0 {
		t.Error("Locations() returned empty slice")
	}

	// Check that all returned locations are valid
	for _, loc := range locs {
		if !loc.Valid() {
			t.Errorf("Locations() returned invalid location: %v", loc)
		}
	}

	// Check that they are sorted by ID (ascending)
	for i := 1; i < len(locs); i++ {
		if locs[i-1].ID() >= locs[i].ID() {
			t.Errorf("Locations() not sorted by ID: %v (ID=%d) should come before %v (ID=%d)",
				locs[i-1], locs[i-1].ID(), locs[i], locs[i].ID())
		}
	}

	// Check that MensaSued (ID=1) is first
	if locs[0] != location.MensaSued {
		t.Errorf("Locations() first element = %v, want %v", locs[0], location.MensaSued)
	}
}

func TestLocation_ID(t *testing.T) {
	tests := []struct {
		name     string
		location location.Location
		want     int
	}{
		{name: "MensaSued", location: location.MensaSued, want: 1},
		{name: "MensaInselschuett", location: location.MensaInselschuett, want: 2},
		{name: "MensaRegensburgerstr", location: location.MensaRegensburgerstr, want: 3},
		{name: "CafeteriaComeIn", location: location.CafeteriaComeIn, want: 10},
		{name: "MensaOic", location: location.MensaOic, want: 28},
		{name: "invalid location", location: location.Location("invalid"), want: 0},
		{name: "empty location", location: location.Location(""), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.location.ID(); got != tt.want {
				t.Errorf("Location.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationOfID(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want location.Location
	}{
		{name: "ID 1 is MensaSued", id: 1, want: location.MensaSued},
		{name: "ID 2 is MensaInselschuett", id: 2, want: location.MensaInselschuett},
		{name: "ID 10 is CafeteriaComeIn", id: 10, want: location.CafeteriaComeIn},
		{name: "ID 28 is MensaOic", id: 28, want: location.MensaOic},
		{name: "ID 0 returns empty", id: 0, want: location.Location("")},
		{name: "ID -1 returns empty", id: -1, want: location.Location("")},
		{name: "ID 999 returns empty", id: 999, want: location.Location("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := location.LocationOfID(tt.id); got != tt.want {
				t.Errorf("LocationOfID(%d) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestLocation_Description(t *testing.T) {
	tests := []struct {
		name        string
		location    location.Location
		wantName    string
		isRefactory bool
		isCafe      bool
		isInternal  bool
	}{
		{
			name:        "MensaSued is refactory",
			location:    location.MensaSued,
			wantName:    "Südmensa",
			isRefactory: true,
			isCafe:      false,
			isInternal:  false,
		},
		{
			name:        "CafeteriaComeIn is cafe",
			location:    location.CafeteriaComeIn,
			wantName:    "Cafeteria \"Come IN\" Hohfederstraße",
			isRefactory: false,
			isCafe:      true,
			isInternal:  false,
		},
		{
			name:        "WohnanlageStPeter is internal",
			location:    location.WohnanlageStPeter,
			wantName:    "Wohnanlage St. Peter",
			isRefactory: false,
			isCafe:      false,
			isInternal:  true,
		},
		{
			name:        "MensaLmp has address",
			location:    location.MensaLmp,
			wantName:    "Mensa Langemarkplatz",
			isRefactory: true,
			isCafe:      false,
			isInternal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := tt.location.Description()
			if desc.Name != tt.wantName {
				t.Errorf("Description().Name = %v, want %v", desc.Name, tt.wantName)
			}
			if desc.Refactory != tt.isRefactory {
				t.Errorf("Description().Refactory = %v, want %v", desc.Refactory, tt.isRefactory)
			}
			if desc.Cafe != tt.isCafe {
				t.Errorf("Description().Cafe = %v, want %v", desc.Cafe, tt.isCafe)
			}
			if desc.Internal != tt.isInternal {
				t.Errorf("Description().Internal = %v, want %v", desc.Internal, tt.isInternal)
			}
		})
	}
}

func TestLocationDescription_Type(t *testing.T) {
	tests := []struct {
		name    string
		desc    location.LocationDescription
		english bool
		want    string
	}{
		{
			name:    "refactory german",
			desc:    location.LocationDescription{Refactory: true},
			english: false,
			want:    "Mensa",
		},
		{
			name:    "refactory english",
			desc:    location.LocationDescription{Refactory: true},
			english: true,
			want:    "Servery",
		},
		{
			name:    "cafe german",
			desc:    location.LocationDescription{Cafe: true},
			english: false,
			want:    "Café",
		},
		{
			name:    "cafe english",
			desc:    location.LocationDescription{Cafe: true},
			english: true,
			want:    "Café",
		},
		{
			name:    "internal german",
			desc:    location.LocationDescription{Internal: true},
			english: false,
			want:    "Intern",
		},
		{
			name:    "internal english",
			desc:    location.LocationDescription{Internal: true},
			english: true,
			want:    "Internal",
		},
		{
			name:    "none returns empty",
			desc:    location.LocationDescription{},
			english: false,
			want:    "",
		},
		{
			name:    "none returns empty english",
			desc:    location.LocationDescription{},
			english: true,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.desc.Type(tt.english); got != tt.want {
				t.Errorf("LocationDescription.Type(%v) = %v, want %v", tt.english, got, tt.want)
			}
		})
	}
}

func TestLocationDescription_Address(t *testing.T) {
	desc := location.LocationDescription{
		Street:   "Erwin-Rommel-Straße",
		StreetNo: "60",
		ZIP:      "91058",
		City:     "Erlangen",
	}

	html := string(desc.Address())

	// Check that it contains the address
	if !strings.Contains(html, "Erwin-Rommel-Straße 60") {
		t.Errorf("Address() does not contain street: %v", html)
	}
	if !strings.Contains(html, "91058 Erlangen") {
		t.Errorf("Address() does not contain city: %v", html)
	}

	// Check that it's a link to OpenStreetMap
	if !strings.Contains(html, "openstreetmap.org/search") {
		t.Errorf("Address() does not contain OSM link: %v", html)
	}

	// Check that it has the proper link attributes
	if !strings.Contains(html, "target='_blank'") {
		t.Errorf("Address() does not open in new tab: %v", html)
	}
	if !strings.Contains(html, "rel='noopener noreferer'") {
		t.Errorf("Address() missing rel attribute: %v", html)
	}
}

func TestLocationDescription_Cmp(t *testing.T) {
	refactory := location.LocationDescription{Name: "A", Refactory: true, City: "A"}
	cafe := location.LocationDescription{Name: "A", Cafe: true, City: "A"}
	internal := location.LocationDescription{Name: "A", Internal: true, City: "A"}

	tests := []struct {
		name string
		a    location.LocationDescription
		b    location.LocationDescription
		want int
	}{
		// Kind comparison
		{name: "refactory < cafe", a: refactory, b: cafe, want: -1},
		{name: "cafe > refactory", a: cafe, b: refactory, want: 1},
		{name: "cafe < internal", a: cafe, b: internal, want: -1},
		{name: "internal > cafe", a: internal, b: cafe, want: 1},
		{name: "refactory < internal", a: refactory, b: internal, want: -1},

		// Same kind, different city
		{
			name: "same kind, city A < city B",
			a:    location.LocationDescription{Refactory: true, City: "Ansbach"},
			b:    location.LocationDescription{Refactory: true, City: "Berlin"},
			want: -1,
		},
		{
			name: "same kind, city B > city A",
			a:    location.LocationDescription{Refactory: true, City: "Berlin"},
			b:    location.LocationDescription{Refactory: true, City: "Ansbach"},
			want: 1,
		},

		// Same kind and city, different name
		{
			name: "same kind and city, name A < name B",
			a:    location.LocationDescription{Refactory: true, City: "Erlangen", Name: "Alpha"},
			b:    location.LocationDescription{Refactory: true, City: "Erlangen", Name: "Beta"},
			want: -1,
		},

		// Equal
		{
			name: "equal",
			a:    location.LocationDescription{Refactory: true, City: "Erlangen", Name: "Test"},
			b:    location.LocationDescription{Refactory: true, City: "Erlangen", Name: "Test"},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Cmp(tt.b); got != tt.want {
				t.Errorf("Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_MarshalJSON(t *testing.T) {
	loc := location.MensaSued
	data, err := json.Marshal(&loc)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Check that ID is present
	if id, ok := result["id"]; !ok || id != "mensa-sued" {
		t.Errorf("MarshalJSON() id = %v, want %v", id, "mensa-sued")
	}

	// Check that Name is present from description
	if name, ok := result["Name"]; !ok || name != "Südmensa" {
		t.Errorf("MarshalJSON() Name = %v, want %v", name, "Südmensa")
	}

	// Check that Refactory is present
	if refactory, ok := result["Refactory"]; !ok || refactory != true {
		t.Errorf("MarshalJSON() Refactory = %v, want %v", refactory, true)
	}
}

func TestLocation_IDRoundTrip(t *testing.T) {
	// Test that LocationOfID(loc.ID()) == loc for all valid locations
	for _, loc := range location.Locations() {
		id := loc.ID()
		if id == 0 {
			t.Errorf("Location %v has ID 0", loc)
			continue
		}
		recovered := location.LocationOfID(id)
		if recovered != loc {
			t.Errorf("LocationOfID(%d) = %v, want %v", id, recovered, loc)
		}
	}
}

func TestAllLocationsHaveDescriptions(t *testing.T) {
	for _, loc := range location.Locations() {
		desc := loc.Description()
		if desc.Name == "" {
			t.Errorf("Location %v has empty Name", loc)
		}
		// Every location should be exactly one type
		typeCount := 0
		if desc.Refactory {
			typeCount++
		}
		if desc.Cafe {
			typeCount++
		}
		if desc.Internal {
			typeCount++
		}
		if typeCount != 1 {
			t.Errorf("Location %v has %d types (should be exactly 1)", loc, typeCount)
		}
	}
}
