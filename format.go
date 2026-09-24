package tidyunits

import (
	"strconv"
	"strings"
)

// formatFloat renders f with up to two decimal places, trimming trailing
// zeros so whole numbers print as "2" instead of "2.00".
func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
