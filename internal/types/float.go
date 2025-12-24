package types

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// SmartFload64 is a data type that can unmarshal from both text and xml elements.
//
// It automatically detects if "," or "." are used as decimal separators.
// The "-" is detected as an empty value, and unmarshaed as 0.
type SmartFloat64 float64

const (
	smartComma  = ","
	smartPeriod = "."
	smartEmpty  = "-"
)

func (sf64 *SmartFloat64) UnmarshalText(text []byte) error {
	value := string(text)

	// if there are only empty values, ignore them
	if value == "" || value == smartEmpty {
		*sf64 = 0
		return nil
	}

	// replace "," by "." (for german notation)
	if strings.ContainsAny(value, smartComma) && !strings.ContainsAny(value, smartPeriod) {
		value = strings.ReplaceAll(value, smartComma, smartPeriod)
	}

	f64, err := strconv.ParseFloat(value, 64)
	*sf64 = SmartFloat64(f64)
	return err
}

func (sf64 *SmartFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	return sf64.UnmarshalText([]byte(attr.Value))
}
