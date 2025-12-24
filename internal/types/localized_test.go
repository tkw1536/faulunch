package types

import "testing"

func TestLPrice_String(t *testing.T) {
	tests := []struct {
		name   string
		price  LPrice
		wantEN string
		wantDE string
	}{
		{name: "zero", price: 0, wantEN: "0.00", wantDE: "0,00"},
		{name: "integer", price: 5, wantEN: "5.00", wantDE: "5,00"},
		{name: "two decimals", price: 3.14, wantEN: "3.14", wantDE: "3,14"},
		{name: "one decimal rounds", price: 2.5, wantEN: "2.50", wantDE: "2,50"},
		{name: "more decimals truncates", price: 1.999, wantEN: "2.00", wantDE: "2,00"},
		{name: "negative", price: -4.50, wantEN: "-4.50", wantDE: "-4,50"},
		{name: "large number", price: 1234.56, wantEN: "1234.56", wantDE: "1234,56"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.ENString(); got != tt.wantEN {
				t.Errorf("LPrice.ENString() = %v, want %v", got, tt.wantEN)
			}
			if got := tt.price.DEString(); got != tt.wantDE {
				t.Errorf("LPrice.DEString() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}

func TestLFloat_String(t *testing.T) {
	tests := []struct {
		name   string
		value  LFloat
		wantEN string
		wantDE string
	}{
		{name: "zero", value: 0, wantEN: "0", wantDE: "0"},
		{name: "integer", value: 5, wantEN: "5", wantDE: "5"},
		{name: "one decimal", value: 2.5, wantEN: "2.5", wantDE: "2,5"},
		{name: "two decimals", value: 3.14, wantEN: "3.14", wantDE: "3,14"},
		{name: "trailing zeros trimmed", value: 1.10, wantEN: "1.1", wantDE: "1,1"},
		{name: "many decimals", value: 1.23456, wantEN: "1.23456", wantDE: "1,23456"},
		{name: "five decimal precision", value: 1.234567, wantEN: "1.23457", wantDE: "1,23457"},
		{name: "negative", value: -3.5, wantEN: "-3.5", wantDE: "-3,5"},
		{name: "large number", value: 1234.5, wantEN: "1234.5", wantDE: "1234,5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.ENString(); got != tt.wantEN {
				t.Errorf("LFloat.ENString() = %v, want %v", got, tt.wantEN)
			}
			if got := tt.value.DEString(); got != tt.wantDE {
				t.Errorf("LFloat.DEString() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}
