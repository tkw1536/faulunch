package types_test

import (
	"encoding/xml"
	"testing"

	"github.com/tkw1536/faulunch/internal/types"
)

func TestSmartFloat64_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    types.SmartFloat64
		wantErr bool
	}{
		// empty values
		{name: "empty string", input: "", want: 0, wantErr: false},
		{name: "dash (empty)", input: "-", want: 0, wantErr: false},

		// period as decimal separator
		{name: "integer", input: "42", want: 42, wantErr: false},
		{name: "float with period", input: "3.14", want: 3.14, wantErr: false},
		{name: "negative float with period", input: "-3.14", want: -3.14, wantErr: false},
		{name: "zero", input: "0", want: 0, wantErr: false},
		{name: "zero with period", input: "0.0", want: 0, wantErr: false},

		// comma as decimal separator (german notation)
		{name: "float with comma", input: "3,14", want: 3.14, wantErr: false},
		{name: "negative float with comma", input: "-3,14", want: -3.14, wantErr: false},
		{name: "zero with comma", input: "0,0", want: 0, wantErr: false},
		{name: "large number with comma", input: "1234,56", want: 1234.56, wantErr: false},

		// mixed (comma and period) - should NOT replace comma
		{name: "mixed comma and period", input: "1.234,56", want: 0, wantErr: true},

		// invalid inputs
		{name: "invalid string", input: "abc", want: 0, wantErr: true},
		{name: "partially invalid", input: "12.34abc", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got types.SmartFloat64
			err := got.UnmarshalText([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSmartFloat64_UnmarshalXMLAttr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    types.SmartFloat64
		wantErr bool
	}{
		{name: "empty attr", input: "", want: 0, wantErr: false},
		{name: "dash attr", input: "-", want: 0, wantErr: false},
		{name: "float with period", input: "2.5", want: 2.5, wantErr: false},
		{name: "float with comma", input: "2,5", want: 2.5, wantErr: false},
		{name: "integer", input: "10", want: 10, wantErr: false},
		{name: "invalid", input: "not a number", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got types.SmartFloat64
			attr := xml.Attr{Value: tt.input}
			err := got.UnmarshalXMLAttr(attr)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXMLAttr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalXMLAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
