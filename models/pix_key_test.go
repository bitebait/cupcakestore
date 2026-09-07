package models

import (
	"strings"
	"testing"
)

func TestNormalizePixKey(t *testing.T) {
	for _, tc := range []struct{ kind, input, want string }{
		{"cpf", "123.456.789-09", "12345678909"},
		{"cnpj", "00.038.166/0001-05", "00038166000105"},
		{"cnpj", "12.abc.345/01de-35", "12ABC34501DE35"},
		{"email", " Loja+Pix@EXAMPLE.COM ", "loja+pix@example.com"},
		{"email", strings.Repeat("a", 64) + "@example.test", strings.Repeat("a", 64) + "@example.test"},
		{"celular", "+55 (61) 91234-5678", "+5561912345678"},
		{"aleatoria", "123E4567-E12B-12D1-A456-426655440000", "123e4567-e12b-12d1-a456-426655440000"},
	} {
		t.Run(tc.kind+"/"+tc.input, func(t *testing.T) {
			got, err := NormalizePixKey(tc.kind, tc.input)
			if err != nil || got != tc.want {
				t.Fatalf("got=%q want=%q err=%v", got, tc.want, err)
			}
		})
	}
}

func TestNormalizePixKeyRejectsInvalidAndMismatchedKeys(t *testing.T) {
	for _, tc := range []struct{ kind, value string }{
		{"cpf", "12345678900"}, {"cpf", "11111111111"}, {"cpf", "loja@example.com"},
		{"cnpj", "12ABC34501DE36"}, {"cnpj", "00000000000000"}, {"cnpj", "12ABC34501DE3A"},
		{"celular", "11999999999"}, {"celular", "+05511999999999"}, {"celular", "+5511<script>"},
		{"aleatoria", "123456"}, {"aleatoria", "00000000-0000-0000-0000-000000000000"},
		{"email", "Nome <loja@example.com>"}, {"email", "loja@example"}, {"email", "joão@example.com"},
		{"email", "a..b@example.com"}, {"email", "a@example.com\nother@example.com"},
		{"email", strings.Repeat("a", 65) + "@example.test"}, {"email", ""}, {"invalid", "loja@example.com"},
	} {
		if value, err := NormalizePixKey(tc.kind, tc.value); err == nil {
			t.Fatalf("accepted %s %q as %q", tc.kind, tc.value, value)
		}
	}
}

func TestNormalizePixMerchantUsesBoundedASCIIWithoutSplittingUnicode(t *testing.T) {
	name, city, err := NormalizePixMerchant("  Doces   São João  ", "São Paulo")
	if err != nil || name != "Doces Sao Joao" || city != "Sao Paulo" {
		t.Fatalf("name=%q city=%q err=%v", name, city, err)
	}
	name, city, err = NormalizePixMerchant(strings.Repeat("Á", 30), strings.Repeat("Ç", 20))
	if err != nil || name != strings.Repeat("A", 25) || city != strings.Repeat("C", 15) {
		t.Fatalf("invalid normalized lengths: %q %q %v", name, city, err)
	}
	for _, pair := range [][2]string{{"", "Sao Paulo"}, {"Loja", " "}, {"Loja🧁", "Sao Paulo"}, {"Loja\nOutra", "Sao Paulo"}, {"\xff", "Sao Paulo"}} {
		if _, _, err := NormalizePixMerchant(pair[0], pair[1]); err == nil {
			t.Fatalf("accepted invalid merchant %q", pair)
		}
	}
}
