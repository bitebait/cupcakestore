package views

import (
	"math"
	"strconv"
	"strings"
)

// Money formats a display amount in Brazilian reais. Form fields keep their
// original numeric values; this helper does not change the stored amount.
func Money(amount float64) string {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return "—"
	}
	formatted := strconv.FormatFloat(amount, 'f', 2, 64)
	negative := strings.HasPrefix(formatted, "-")
	parts := strings.SplitN(strings.TrimPrefix(formatted, "-"), ".", 2)
	var whole strings.Builder
	for i, digit := range parts[0] {
		if i > 0 && (len(parts[0])-i)%3 == 0 {
			whole.WriteByte('.')
		}
		whole.WriteRune(digit)
	}
	prefix := "R$\u00a0"
	if negative && (parts[0] != "0" || parts[1] != "00") {
		prefix = "-" + prefix
	}
	return prefix + whole.String() + "," + parts[1]
}
