package ltime

//spellchecker:words database driver html template strconv time gorm schema tzdata
import (
	"database/sql/driver"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "time/tzdata"
)

// Day represents a single day, represented as a unix timestamp
type Day int64

// ParseDay parses a day from an unknown input value.
//
// The input value may be a numeric type, a string, a []byte, or a time.Time.
// If the underlying value cannot be parsed, the zero Day is returned.
func ParseDay(value any) Day {
	var di int64
	var err error
	switch v := value.(type) {
	default:
		// di = 0 (unknown type)
	case time.Time:
		di = v.Unix()

	// string, byte
	case string:
		di, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			di = 0
		}
	case []byte:
		di, err = strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			di = 0
		}
	// ints
	case int:
		di = int64(v)
	case int8:
		di = int64(v)
	case int16:
		di = int64(v)
	case int32:
		di = int64(v)
	case int64:
		di = v
	// uints
	case uint:
		di = int64(v)
	case uint8:
		di = int64(v)
	case uint16:
		di = int64(v)
	case uint32:
		di = int64(v)
	case uint64:
		di = int64(v)
	}

	// ensure that the day is positive
	if di < 0 {
		di = 0
	}

	// and return the day
	return Day(di)
}

// Today returns the current day
func Today() Day {
	return normalizeDay(time.Now())
}

var europeBerlin *time.Location

func init() {
	var err error
	europeBerlin, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
}

func normalizeDay(t time.Time) Day {
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, europeBerlin)
	return Day(d.Unix())
}

// Normalize normalizes the given day, that is sets it to 0:00 of the local timezone.
func (d Day) Normalize() Day {
	return normalizeDay(d.Time())
}

// Add adds count number of days to the current day.
// The result is normalized.
func (d Day) Add(count int) Day {
	t := d.Time().Add(time.Duration(count) * 24 * time.Hour)
	return normalizeDay(t)
}

// Time returns the time behind this day in the appropriate local timezone.
func (d Day) Time() time.Time {
	// TODO: Do we need this?
	return time.Unix(int64(d), 0).In(europeBerlin)
}

func (d *Day) Scan(value any) error {
	*d = ParseDay(value)
	return nil
}

func (d Day) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "integer"
}

// Equal compares if this day is equal to another day
func (d Day) Equal(other Day) bool {
	return int64(d) == int64(other)
}

// Day returns driver value, being int64
func (d Day) Value() (driver.Value, error) {
	return int64(d), nil
}

var daysDE = [...]string{
	"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag",
}
var monthsDE = [...]string{
	"Januar", "Februar", "März", "April", "Mai", "Juni",
	"Juli", "August", "September", "Oktober", "November", "Dezember",
}

// String formats this Day as a number of seconds since the unix epoch.
func (d Day) String() string {
	return strconv.FormatInt(int64(d), 10)
}

// DEString formats this date as a german localized string.
func (d Day) DEString() string {
	t := d.Time()
	return fmt.Sprintf("%s, %d. %s %04d",
		daysDE[t.Weekday()], t.Day(), monthsDE[t.Month()-1], t.Year(),
	)
}

const dateStamp = "2006-01-02"

func (d Day) DEHTML() template.HTML {
	return template.HTML("<time datetime='" + d.Time().Format(dateStamp) + "'>" + d.DEString() + "</time>")
}

var daysEN = [...]string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}
var monthsEN = [...]string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// ENString formats this day as an english localized string.
func (d Day) ENString() string {
	t := d.Time()

	suffix := "th"
	day := t.Day()
	if day == 1 || day == 21 || day == 31 {
		suffix = "st"
	}
	if day == 2 || day == 22 {
		suffix = "nd"
	}
	if day == 3 || day == 23 {
		suffix = "rd"
	}
	return fmt.Sprintf("%s, %d%s %s %04d",
		daysEN[t.Weekday()], t.Day(), suffix, monthsEN[t.Month()-1], t.Year(),
	)
}

func (d Day) ENHTML() template.HTML {
	return template.HTML("<time datetime='" + d.Time().Format(dateStamp) + "'>" + d.ENString() + "</time>")
}
