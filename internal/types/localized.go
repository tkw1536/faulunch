package types

import (
	"strconv"
	"strings"
)

// LPrice represents a localizable price.
type LPrice float64

func (lp LPrice) DEString() string {
	return strings.ReplaceAll(lp.ENString(), ".", ",")
}
func (lp LPrice) ENString() string {
	return strconv.FormatFloat(float64(lp), 'f', 2, 64)
}

// LFloat represents a localized float
type LFloat float64

func (lf LFloat) DEString() string {
	return strings.ReplaceAll(lf.ENString(), ".", ",")
}

func (lf LFloat) ENString() string {
	value := strconv.FormatFloat(float64(lf), 'f', 5, 64)
	return strings.TrimSuffix(strings.TrimRight(value, "0"), ".")
}
