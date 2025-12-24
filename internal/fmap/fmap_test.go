package fmap_test

import (
	"testing"

	"github.com/tkw1536/faulunch/internal/fmap"
)

func TestFMap_Add(t *testing.T) {
	// Note: Add() only updates the value if the key already exists (case-insensitive).
	// It returns true if the key was found and updated, false if the key was not found.
	tests := []struct {
		name      string
		initial   fmap.FMap[string, int]
		key       string
		value     int
		wantNew   bool
		wantValue int
	}{
		{
			name:      "add to empty map - key not found",
			initial:   fmap.FMap[string, int]{},
			key:       "Hello",
			value:     1,
			wantNew:   false,
			wantValue: 0, // not added
		},
		{
			name:      "update existing key same case",
			initial:   fmap.FMap[string, int]{"Hello": 1},
			key:       "Hello",
			value:     2,
			wantNew:   true,
			wantValue: 2,
		},
		{
			name:      "update existing key different case",
			initial:   fmap.FMap[string, int]{"Hello": 1},
			key:       "HELLO",
			value:     3,
			wantNew:   true,
			wantValue: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.initial.Add(tt.key, tt.value)
			if got != tt.wantNew {
				t.Errorf("Add() = %v, want %v", got, tt.wantNew)
			}
			_, val, ok := tt.initial.Get(tt.key)
			if tt.wantValue == 0 && !ok {
				// expected: key not found
				return
			}
			if val != tt.wantValue {
				t.Errorf("Get(%q) value = %v, want %v", tt.key, val, tt.wantValue)
			}
		})
	}
}

func TestFMap_Remove(t *testing.T) {
	tests := []struct {
		name    string
		initial fmap.FMap[string, int]
		key     string
		wantOk  bool
		wantHas bool
	}{
		{
			name:    "remove existing key same case",
			initial: fmap.FMap[string, int]{"Hello": 1},
			key:     "Hello",
			wantOk:  true,
			wantHas: false,
		},
		{
			name:    "remove existing key different case",
			initial: fmap.FMap[string, int]{"Hello": 1},
			key:     "HELLO",
			wantOk:  true,
			wantHas: false,
		},
		{
			name:    "remove non-existing key",
			initial: fmap.FMap[string, int]{"Hello": 1},
			key:     "World",
			wantOk:  false,
			wantHas: true,
		},
		{
			name:    "remove from empty map",
			initial: fmap.FMap[string, int]{},
			key:     "Hello",
			wantOk:  false,
			wantHas: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.initial.Remove(tt.key)
			if got != tt.wantOk {
				t.Errorf("Remove() = %v, want %v", got, tt.wantOk)
			}
			if has := tt.initial.Has("Hello"); has != tt.wantHas {
				t.Errorf("Has(Hello) after Remove = %v, want %v", has, tt.wantHas)
			}
		})
	}
}

func TestFMap_Has(t *testing.T) {
	fm := fmap.FMap[string, int]{"Hello": 1, "World": 2}

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "exact match", key: "Hello", want: true},
		{name: "lowercase", key: "hello", want: true},
		{name: "uppercase", key: "HELLO", want: true},
		{name: "mixed case", key: "hElLo", want: true},
		{name: "not found", key: "Foo", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fm.Has(tt.key); got != tt.want {
				t.Errorf("Has(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestFMap_Key(t *testing.T) {
	fm := fmap.FMap[string, int]{"Hello": 1}

	tests := []struct {
		name    string
		key     string
		wantKey string
		wantOk  bool
	}{
		{name: "exact match", key: "Hello", wantKey: "Hello", wantOk: true},
		{name: "lowercase", key: "hello", wantKey: "Hello", wantOk: true},
		{name: "uppercase", key: "HELLO", wantKey: "Hello", wantOk: true},
		{name: "not found returns input", key: "World", wantKey: "World", wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotOk := fm.Key(tt.key)
			if gotKey != tt.wantKey {
				t.Errorf("Key(%q) key = %v, want %v", tt.key, gotKey, tt.wantKey)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Key(%q) ok = %v, want %v", tt.key, gotOk, tt.wantOk)
			}
		})
	}
}

func TestFMap_Get(t *testing.T) {
	fm := fmap.FMap[string, int]{"Hello": 42}

	tests := []struct {
		name      string
		key       string
		wantKey   string
		wantValue int
		wantOk    bool
	}{
		{name: "exact match", key: "Hello", wantKey: "Hello", wantValue: 42, wantOk: true},
		{name: "lowercase", key: "hello", wantKey: "Hello", wantValue: 42, wantOk: true},
		{name: "uppercase", key: "HELLO", wantKey: "Hello", wantValue: 42, wantOk: true},
		{name: "not found", key: "World", wantKey: "World", wantValue: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, gotOk := fm.Get(tt.key)
			if gotKey != tt.wantKey {
				t.Errorf("Get(%q) key = %v, want %v", tt.key, gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("Get(%q) value = %v, want %v", tt.key, gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Get(%q) ok = %v, want %v", tt.key, gotOk, tt.wantOk)
			}
		})
	}
}
