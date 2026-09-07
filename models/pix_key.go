package models

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

var (
	pixCPFPattern    = regexp.MustCompile(`^[0-9]{11}$`)
	pixCNPJPattern   = regexp.MustCompile(`^[A-Z0-9]{12}[0-9]{2}$`)
	pixPhonePattern  = regexp.MustCompile(`^\+[1-9][0-9]{1,14}$`)
	pixEVPPattern    = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	pixDomainPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)+$`)
)

// NormalizePixKey validates the key's format/check digits, not its registration
// or ownership in DICT. The merchant must verify the destination in their bank.
func NormalizePixKey(kind, key string) (string, error) {
	key = strings.TrimSpace(key)
	if !utf8.ValidString(key) || key == "" {
		return "", errors.New("informe uma chave Pix válida")
	}
	var valid bool
	switch pixType(kind) {
	case PixTypeCPF:
		key = strings.NewReplacer(".", "", "-", "").Replace(key)
		valid = pixCPFPattern.MatchString(key) && validPixCPF(key)
	case PixTypeCNPJ:
		key = strings.ToUpper(strings.NewReplacer(".", "", "/", "", "-", "").Replace(key))
		valid = pixCNPJPattern.MatchString(key) && validPixCNPJ(key)
	case PixTypePhone:
		key = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "").Replace(key)
		if !pixPhonePattern.MatchString(key) {
			return "", errors.New("informe o celular Pix com código do país, por exemplo +5511999999999")
		}
		valid = true
	case PixTypeRandomKey:
		key = strings.ToLower(key)
		valid = pixEVPPattern.MatchString(key) && key != "00000000-0000-0000-0000-000000000000"
	case PixTypeEmail:
		key = strings.ToLower(key)
		address, err := mail.ParseAddress(key)
		local, domain, found := strings.Cut(key, "@")
		valid = err == nil && address.Address == key && found && len(local) <= 64 &&
			!strings.ContainsAny(local, "\"\\") && !strings.Contains(local, "..") &&
			!strings.HasPrefix(local, ".") && !strings.HasSuffix(local, ".") && pixDomainPattern.MatchString(domain)
	default:
		return "", errors.New("selecione um tipo de chave Pix válido")
	}
	if !valid || len(key) > 77 || !pixASCII(key) {
		return "", errors.New("a chave Pix não corresponde ao tipo selecionado ou excede o limite de 77 caracteres")
	}
	return key, nil
}

func validPixCPF(value string) bool {
	if value == strings.Repeat(value[:1], 11) {
		return false
	}
	for length := 9; length <= 10; length++ {
		sum := 0
		for i := 0; i < length; i++ {
			sum += int(value[i]-'0') * (length + 1 - i)
		}
		if int(value[length]-'0') != pixCheckDigit(sum) {
			return false
		}
	}
	return true
}

// Receita Federal's alphanumeric CNPJ algorithm maps each character to ASCII-48.
// https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/documentos-tecnicos/cnpj/manual-dv-cnpj.pdf
func validPixCNPJ(value string) bool {
	if value == strings.Repeat(value[:1], 14) {
		return false
	}
	for length := 12; length <= 13; length++ {
		sum, weight := 0, 2
		for i := length - 1; i >= 0; i-- {
			sum += int(value[i]-'0') * weight
			weight++
			if weight > 9 {
				weight = 2
			}
		}
		if int(value[length]-'0') != pixCheckDigit(sum) {
			return false
		}
	}
	return true
}

func pixCheckDigit(sum int) int {
	if remainder := sum % 11; remainder >= 2 {
		return 11 - remainder
	}
	return 0
}

// NormalizePixMerchant maps Portuguese diacritics to the ASCII BR Code fields.
// Truncation happens after normalization, so no UTF-8 sequence can be split.
func NormalizePixMerchant(name, city string) (string, string, error) {
	name, err := normalizePixText(name, 25)
	if err != nil {
		return "", "", errors.New("informe o nome do recebedor Pix usando letras, números e pontuação simples")
	}
	city, err = normalizePixText(city, 15)
	if err != nil {
		return "", "", errors.New("informe a cidade do recebedor Pix usando letras e pontuação simples")
	}
	return name, city, nil
}

func normalizePixText(value string, limit int) (string, error) {
	if !utf8.ValidString(value) {
		return "", errors.New("texto inválido")
	}
	var result strings.Builder
	for _, character := range norm.NFD.String(value) {
		if unicode.Is(unicode.Mn, character) {
			continue
		}
		if character < 32 || character > 126 {
			return "", errors.New("caractere não suportado")
		}
		result.WriteRune(character)
	}
	text := strings.Join(strings.Fields(result.String()), " ")
	if len(text) > limit {
		text = strings.TrimSpace(text[:limit])
	}
	if text == "" {
		return "", errors.New("texto obrigatório")
	}
	return text, nil
}

func pixASCII(value string) bool {
	for _, character := range value {
		if character < 33 || character > 126 {
			return false
		}
	}
	return true
}
