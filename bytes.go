package tidyunits

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// byteUnitMultipliers maps a normalized unit token to the number of bytes
// it represents. Both decimal (kb, mb, ...) and binary (kib, mib, ...)
// units are accepted on input, since real-world data mixes both freely.
var byteUnitMultipliers = map[string]float64{
	"":      1,
	"b":     1,
	"byte":  1,
	"bytes": 1,

	"kb": 1000,
	"mb": 1000 * 1000,
	"gb": 1000 * 1000 * 1000,
	"tb": 1000 * 1000 * 1000 * 1000,
	"pb": 1000 * 1000 * 1000 * 1000 * 1000,

	"kib": 1024,
	"mib": 1024 * 1024,
	"gib": 1024 * 1024 * 1024,
	"tib": 1024 * 1024 * 1024 * 1024,
	"pib": 1024 * 1024 * 1024 * 1024 * 1024,
}

// decimalByteSteps drives FormatBytes. Only decimal units are produced on
// output, even though ParseBytes accepts binary ones on input: it's the
// convention most readers expect from a plain byte count.
var decimalByteSteps = []struct {
	suffix string
	size   float64
}{
	{"PB", 1000 * 1000 * 1000 * 1000 * 1000},
	{"TB", 1000 * 1000 * 1000 * 1000},
	{"GB", 1000 * 1000 * 1000},
	{"MB", 1000 * 1000},
	{"kB", 1000},
}

// ParseBytes turns a human-written byte size into an exact byte count.
// It accepts optional whitespace between the number and unit, is
// case-insensitive, and understands both decimal units (kb, mb, gb, ...)
// and binary units (kib, mib, gib, ...). A bare number is treated as a
// byte count. Examples: "1.5GB", "2048", "3 MiB", "700kb".
func ParseBytes(s string) (int64, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, fmt.Errorf("tidyunits: empty byte size")
	}

	i := 0
	for i < len(trimmed) && (trimmed[i] == '.' || trimmed[i] == '-' || trimmed[i] == '+' || (trimmed[i] >= '0' && trimmed[i] <= '9')) {
		i++
	}
	if i == 0 {
		return 0, fmt.Errorf("tidyunits: %q has no numeric value", s)
	}

	numPart := trimmed[:i]
	unitPart := strings.ToLower(strings.TrimSpace(trimmed[i:]))

	value, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, fmt.Errorf("tidyunits: %q is not a valid number: %w", numPart, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("tidyunits: byte size %q must not be negative", s)
	}

	multiplier, ok := byteUnitMultipliers[unitPart]
	if !ok {
		return 0, fmt.Errorf("tidyunits: unrecognized byte unit %q", unitPart)
	}

	return int64(math.Round(value * multiplier)), nil
}

// FormatBytes renders a byte count using the largest decimal unit that
// keeps it at least 1, e.g. 1500000 becomes "1.5 MB". Values under 1000
// are rendered as a plain byte count.
func FormatBytes(n int64) string {
	if n < 0 {
		return "-" + FormatBytes(-n)
	}

	f := float64(n)
	for _, step := range decimalByteSteps {
		if f >= step.size {
			return formatFloat(f/step.size) + " " + step.suffix
		}
	}

	if n == 1 {
		return "1 byte"
	}
	return strconv.FormatInt(n, 10) + " bytes"
}
