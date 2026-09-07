package views

import (
	"math"
	"testing"
)

func TestMoney(t *testing.T) {
	for _, test := range []struct {
		name   string
		amount float64
		want   string
	}{
		{"zero", 0, "R$\u00a00,00"},
		{"fraction", 12.5, "R$\u00a012,50"},
		{"cents", 0.07, "R$\u00a00,07"},
		{"rounding", 19.999, "R$\u00a020,00"},
		{"thousands", 1234.56, "R$\u00a01.234,56"},
		{"millions", 1234567.89, "R$\u00a01.234.567,89"},
		{"negative", -1234.56, "-R$\u00a01.234,56"},
		{"negative zero", math.Copysign(0, -1), "R$\u00a00,00"},
		{"nan", math.NaN(), "—"},
		{"infinite", math.Inf(1), "—"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := Money(test.amount); got != test.want {
				t.Errorf("Money(%v) = %q; want %q", test.amount, got, test.want)
			}
		})
	}
}
