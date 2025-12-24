package ltime_test

import (
	"testing"
	"time"

	"github.com/tkw1536/faulunch/internal/ltime"
)

func TestParseDay(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  ltime.Day
	}{
		// time.Time
		{name: "time.Time", input: time.Unix(1609459200, 0), want: ltime.Day(1609459200)},

		// string
		{name: "string valid", input: "1609459200", want: ltime.Day(1609459200)},
		{name: "string invalid", input: "not a number", want: ltime.Day(0)},
		{name: "string empty", input: "", want: ltime.Day(0)},

		// []byte
		{name: "[]byte valid", input: []byte("1609459200"), want: ltime.Day(1609459200)},
		{name: "[]byte invalid", input: []byte("invalid"), want: ltime.Day(0)},

		// int types
		{name: "int", input: int(1609459200), want: ltime.Day(1609459200)},
		{name: "int8", input: int8(100), want: ltime.Day(100)},
		{name: "int16", input: int16(1000), want: ltime.Day(1000)},
		{name: "int32", input: int32(1609459200), want: ltime.Day(1609459200)},
		{name: "int64", input: int64(1609459200), want: ltime.Day(1609459200)},

		// uint types
		{name: "uint", input: uint(1609459200), want: ltime.Day(1609459200)},
		{name: "uint8", input: uint8(100), want: ltime.Day(100)},
		{name: "uint16", input: uint16(1000), want: ltime.Day(1000)},
		{name: "uint32", input: uint32(1609459200), want: ltime.Day(1609459200)},
		{name: "uint64", input: uint64(1609459200), want: ltime.Day(1609459200)},

		// negative becomes zero
		{name: "negative int", input: int(-100), want: ltime.Day(0)},
		{name: "negative int64", input: int64(-1609459200), want: ltime.Day(0)},

		// unknown type
		{name: "unknown type", input: struct{}{}, want: ltime.Day(0)},
		{name: "nil", input: nil, want: ltime.Day(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ltime.ParseDay(tt.input); got != tt.want {
				t.Errorf("ParseDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay_Normalize(t *testing.T) {
	// Create a day from a specific timestamp (2021-01-01 12:30:00 UTC)
	d := ltime.Day(1609504200)
	normalized := d.Normalize()

	// The normalized day should be at midnight Berlin time
	normalizedTime := normalized.Time()
	if normalizedTime.Hour() != 0 || normalizedTime.Minute() != 0 || normalizedTime.Second() != 0 {
		t.Errorf("Normalize() did not set time to midnight, got %v", normalizedTime)
	}
}

func TestDay_Add(t *testing.T) {
	// Use a known date: 2021-01-01 00:00:00 Berlin time
	berlin, _ := time.LoadLocation("Europe/Berlin")
	baseTime := time.Date(2021, 1, 1, 0, 0, 0, 0, berlin)
	d := ltime.Day(baseTime.Unix())

	tests := []struct {
		name  string
		count int
		want  time.Time
	}{
		{name: "add 0 days", count: 0, want: time.Date(2021, 1, 1, 0, 0, 0, 0, berlin)},
		{name: "add 1 day", count: 1, want: time.Date(2021, 1, 2, 0, 0, 0, 0, berlin)},
		{name: "add 7 days", count: 7, want: time.Date(2021, 1, 8, 0, 0, 0, 0, berlin)},
		{name: "subtract 1 day", count: -1, want: time.Date(2020, 12, 31, 0, 0, 0, 0, berlin)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Add(tt.count)
			resultTime := result.Time()
			if !resultTime.Equal(tt.want) {
				t.Errorf("Add(%d) = %v, want %v", tt.count, resultTime, tt.want)
			}
		})
	}
}

func TestDay_Equal(t *testing.T) {
	d1 := ltime.Day(1609459200)
	d2 := ltime.Day(1609459200)
	d3 := ltime.Day(1609545600)

	if !d1.Equal(d2) {
		t.Errorf("Equal() should return true for same values")
	}
	if d1.Equal(d3) {
		t.Errorf("Equal() should return false for different values")
	}
}

func TestDay_Value(t *testing.T) {
	d := ltime.Day(1609459200)
	val, err := d.Value()
	if err != nil {
		t.Errorf("Value() returned error: %v", err)
	}
	if val != int64(1609459200) {
		t.Errorf("Value() = %v, want %v", val, int64(1609459200))
	}
}

func TestDay_Scan(t *testing.T) {
	var d ltime.Day
	err := d.Scan(int64(1609459200))
	if err != nil {
		t.Errorf("Scan() returned error: %v", err)
	}
	if d != ltime.Day(1609459200) {
		t.Errorf("Scan() = %v, want %v", d, ltime.Day(1609459200))
	}
}

func TestDay_String(t *testing.T) {
	tests := []struct {
		name string
		day  ltime.Day
		want string
	}{
		{name: "zero", day: ltime.Day(0), want: "0"},
		{name: "positive", day: ltime.Day(1609459200), want: "1609459200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.day.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay_LocalizedString(t *testing.T) {
	berlin, _ := time.LoadLocation("Europe/Berlin")

	tests := []struct {
		name   string
		day    ltime.Day
		wantEN string
		wantDE string
	}{
		{
			name:   "Jan 1st 2021 (Friday)",
			day:    ltime.Day(time.Date(2021, 1, 1, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Friday, 1st January 2021",
			wantDE: "Freitag, 1. Januar 2021",
		},
		{
			name:   "Feb 2nd 2021 (Tuesday)",
			day:    ltime.Day(time.Date(2021, 2, 2, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Tuesday, 2nd February 2021",
			wantDE: "Dienstag, 2. Februar 2021",
		},
		{
			name:   "Mar 3rd 2021 (Wednesday)",
			day:    ltime.Day(time.Date(2021, 3, 3, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Wednesday, 3rd March 2021",
			wantDE: "Mittwoch, 3. März 2021",
		},
		{
			name:   "Apr 4th 2021 (Sunday)",
			day:    ltime.Day(time.Date(2021, 4, 4, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Sunday, 4th April 2021",
			wantDE: "Sonntag, 4. April 2021",
		},
		{
			name:   "May 21st 2021 (Friday)",
			day:    ltime.Day(time.Date(2021, 5, 21, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Friday, 21st May 2021",
			wantDE: "Freitag, 21. Mai 2021",
		},
		{
			name:   "Jun 22nd 2021 (Tuesday)",
			day:    ltime.Day(time.Date(2021, 6, 22, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Tuesday, 22nd June 2021",
			wantDE: "Dienstag, 22. Juni 2021",
		},
		{
			name:   "Jul 23rd 2021 (Friday)",
			day:    ltime.Day(time.Date(2021, 7, 23, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Friday, 23rd July 2021",
			wantDE: "Freitag, 23. Juli 2021",
		},
		{
			name:   "Aug 31st 2021 (Tuesday)",
			day:    ltime.Day(time.Date(2021, 8, 31, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Tuesday, 31st August 2021",
			wantDE: "Dienstag, 31. August 2021",
		},
		{
			name:   "Sep 15th 2021 (Wednesday)",
			day:    ltime.Day(time.Date(2021, 9, 15, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Wednesday, 15th September 2021",
			wantDE: "Mittwoch, 15. September 2021",
		},
		{
			name:   "Oct 11th 2021 (Monday)",
			day:    ltime.Day(time.Date(2021, 10, 11, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Monday, 11th October 2021",
			wantDE: "Montag, 11. Oktober 2021",
		},
		{
			name:   "Nov 12th 2021 (Friday)",
			day:    ltime.Day(time.Date(2021, 11, 12, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Friday, 12th November 2021",
			wantDE: "Freitag, 12. November 2021",
		},
		{
			name:   "Dec 25th 2021 (Saturday)",
			day:    ltime.Day(time.Date(2021, 12, 25, 0, 0, 0, 0, berlin).Unix()),
			wantEN: "Saturday, 25th December 2021",
			wantDE: "Samstag, 25. Dezember 2021",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.day.ENString(); got != tt.wantEN {
				t.Errorf("ENString() = %v, want %v", got, tt.wantEN)
			}
			if got := tt.day.DEString(); got != tt.wantDE {
				t.Errorf("DEString() = %v, want %v", got, tt.wantDE)
			}
		})
	}
}

func TestDay_LocalizedHTML(t *testing.T) {
	berlin, _ := time.LoadLocation("Europe/Berlin")
	d := ltime.Day(time.Date(2021, 1, 1, 0, 0, 0, 0, berlin).Unix())

	enHTML := d.ENHTML()
	if string(enHTML) != "<time datetime='2021-01-01'>Friday, 1st January 2021</time>" {
		t.Errorf("ENHTML() = %v", enHTML)
	}

	deHTML := d.DEHTML()
	if string(deHTML) != "<time datetime='2021-01-01'>Freitag, 1. Januar 2021</time>" {
		t.Errorf("DEHTML() = %v", deHTML)
	}
}
