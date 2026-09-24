package tidyunits

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// durationUnitMultipliers maps every unit spelling ParseDuration accepts
// to the time.Duration it represents. "d"/"day" isn't a real time.Duration
// constant, so it's expressed as a fixed 24h - fine for messy human input,
// wrong for calendar arithmetic across DST changes, which this package
// never attempts.
var durationUnitMultipliers = map[string]time.Duration{
	"ns": time.Nanosecond,
	"us": time.Microsecond,
	"µs": time.Microsecond,
	"ms": time.Millisecond,

	"s":       time.Second,
	"sec":     time.Second,
	"secs":    time.Second,
	"second":  time.Second,
	"seconds": time.Second,

	"m":       time.Minute,
	"min":     time.Minute,
	"mins":    time.Minute,
	"minute":  time.Minute,
	"minutes": time.Minute,

	"h":     time.Hour,
	"hr":    time.Hour,
	"hrs":   time.Hour,
	"hour":  time.Hour,
	"hours": time.Hour,

	"d":    24 * time.Hour,
	"day":  24 * time.Hour,
	"days": 24 * time.Hour,
}

// ParseDuration turns a human-written duration into a time.Duration. It
// accepts one or more number+unit pairs separated by optional whitespace
// or commas, e.g. "1h30m", "2 hours 15 mins", "90s", "1.5 days". Units are
// case-insensitive. A bare number with no unit is an error: there's no
// safe default to guess.
func ParseDuration(s string) (time.Duration, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, fmt.Errorf("tidyunits: empty duration")
	}

	var total time.Duration
	found := false

	rest := trimmed
	for {
		rest = strings.TrimLeft(rest, " \t,")
		if rest == "" {
			break
		}

		i := 0
		for i < len(rest) && (rest[i] == '.' || (rest[i] >= '0' && rest[i] <= '9')) {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("tidyunits: %q has a value with no leading number", s)
		}
		numPart := rest[:i]
		rest = strings.TrimLeft(rest[i:], " \t")

		j := 0
		for j < len(rest) && !(rest[j] >= '0' && rest[j] <= '9') && rest[j] != ' ' && rest[j] != '\t' && rest[j] != ',' && rest[j] != '.' {
			j++
		}
		if j == 0 {
			return 0, fmt.Errorf("tidyunits: %q is missing a unit after %q", s, numPart)
		}
		unitPart := strings.ToLower(rest[:j])
		rest = rest[j:]

		value, err := strconv.ParseFloat(numPart, 64)
		if err != nil {
			return 0, fmt.Errorf("tidyunits: %q is not a valid number: %w", numPart, err)
		}

		unit, ok := durationUnitMultipliers[unitPart]
		if !ok {
			return 0, fmt.Errorf("tidyunits: unrecognized duration unit %q", unitPart)
		}

		total += time.Duration(value * float64(unit))
		found = true
	}

	if !found {
		return 0, fmt.Errorf("tidyunits: %q has no number+unit pairs", s)
	}

	return total, nil
}

// FormatDuration renders a duration using the single largest whole unit
// that keeps it readable, e.g. 90*time.Minute becomes "1.5h" and
// 500*time.Millisecond becomes "500ms". This differs deliberately from
// time.Duration.String, which always mixes units (e.g. "1h30m0s").
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + FormatDuration(-d)
	}

	switch {
	case d < time.Microsecond:
		return strconv.FormatInt(int64(d), 10) + "ns"
	case d < time.Millisecond:
		return formatFloat(float64(d)/float64(time.Microsecond)) + "µs"
	case d < time.Second:
		return formatFloat(float64(d)/float64(time.Millisecond)) + "ms"
	case d < time.Minute:
		return formatFloat(float64(d)/float64(time.Second)) + "s"
	case d < time.Hour:
		return formatFloat(float64(d)/float64(time.Minute)) + "m"
	case d < 24*time.Hour:
		return formatFloat(float64(d)/float64(time.Hour)) + "h"
	default:
		return formatFloat(float64(d)/float64(24*time.Hour)) + "d"
	}
}
