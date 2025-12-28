//spellchecker:words faulunch
package plan_test

//spellchecker:words bytes http reflect testing github faulunch internal location plan
import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/tkw1536/faulunch/internal/location"
	"github.com/tkw1536/faulunch/internal/plan"
)

// mockClient is a mock implementation of plan.ClientLike for testing.
type mockClient struct {
	response *http.Response
	err      error
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestPlanURL(t *testing.T) {
	tests := []struct {
		name    string
		loc     location.Location
		english bool
		want    string
	}{
		{
			name:    "MensaSued German",
			loc:     location.MensaSued,
			english: false,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/mensa-sued.xml",
		},
		{
			name:    "MensaSued English",
			loc:     location.MensaSued,
			english: true,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/en/mensa-sued.xml",
		},
		{
			name:    "MensaInselschuett German",
			loc:     location.MensaInselschuett,
			english: false,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/mensa-inselschuett.xml",
		},
		{
			name:    "MensaInselschuett English",
			loc:     location.MensaInselschuett,
			english: true,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/en/mensa-inselschuett.xml",
		},
		{
			name:    "CafeteriaComeIn German",
			loc:     location.CafeteriaComeIn,
			english: false,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/cafeteria-come-in.xml",
		},
		{
			name:    "CafeteriaComeIn English",
			loc:     location.CafeteriaComeIn,
			english: true,
			want:    "https://www.max-manager.de/daten-extern/sw-erlangen-nuernberg/xml/en/cafeteria-come-in.xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := plan.PlanURL(tt.loc, tt.english); got != tt.want {
				t.Errorf("PlanURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	validXML := `<?xml version='1.0' encoding='utf-8'?>
<speiseplan locationId='42'>
<tag timestamp='1234567890'>
<item>
<category>Test Category</category>
<title>Test Dish (A,B,C)</title>
<description>A test description</description>
<beilagen>Test side dish</beilagen>
<preis1>1,50</preis1>
<preis2>2,50</preis2>
<preis3>3,50</preis3>
<einheit>portion</einheit>
<piktogramme>test-icon</piktogramme>
<kj>100.0</kj>
<kcal>25.0</kcal>
<fett>1.5</fett>
<gesfett>0.5</gesfett>
<kh>3.0</kh>
<zucker>1.0</zucker>
<ballaststoffe>0.2</ballaststoffe>
<eiweiss>2.0</eiweiss>
<salz>0.1</salz>
<foto>test.jpg</foto>
</item>
</tag>
</speiseplan>`

	// Parse the expected plan from the same XML
	var wantPlan plan.Plan
	if err := xml.Unmarshal([]byte(validXML), &wantPlan); err != nil {
		t.Fatalf("failed to parse expected plan: %v", err)
	}

	tests := []struct {
		name    string
		client  *mockClient
		loc     location.Location
		english bool
		wantErr bool
		want    plan.Plan
	}{
		{
			name: "successful fetch",
			client: &mockClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(validXML)),
				},
			},
			loc:     location.MensateriaOhm,
			english: false,
			wantErr: false,
			want:    wantPlan,
		},
		{
			name: "non-OK status code",
			client: &mockClient{
				response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
			},
			loc:     location.MensateriaOhm,
			english: false,
			wantErr: true,
		},
		{
			name: "invalid XML",
			client: &mockClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("<invalid>")),
				},
			},
			loc:     location.MensateriaOhm,
			english: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := plan.Fetch(t.Context(), tt.client, tt.loc, tt.english)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
