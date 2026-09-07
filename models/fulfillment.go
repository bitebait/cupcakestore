package models

import (
	"errors"
	"strings"
)

var (
	ErrContactIncomplete         = errors.New("complete seu nome, sobrenome e telefone para continuar")
	ErrDeliveryAddressIncomplete = errors.New("complete seu endereço, cidade, UF e CEP para receber a entrega")
	ErrDeliveryUnavailable       = errors.New("a entrega está indisponível; revise a forma de recebimento")
	ErrPickupUnavailable         = errors.New("a retirada está indisponível: a loja precisa informar endereço, cidade e UF válidos")
	ErrFulfillmentUnavailable    = errors.New("a loja não possui uma forma de recebimento disponível no momento")
	ErrPaymentUnavailable        = errors.New("a loja não possui uma forma de pagamento disponível no momento")
)

func validBrazilianState(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO":
		return true
	}
	return false
}

func validAddress(address, city, state string) bool {
	return strings.TrimSpace(address) != "" && strings.TrimSpace(city) != "" && validBrazilianState(state)
}

func formattedDigits(value string) string {
	var digits strings.Builder
	for _, char := range strings.TrimSpace(value) {
		switch {
		case char >= '0' && char <= '9':
			digits.WriteRune(char)
		case char == ' ' || char == '-' || char == '(' || char == ')' || char == '+':
		default:
			return ""
		}
	}
	return digits.String()
}

func validPhone(number string) bool {
	digits := formattedDigits(number)
	if (len(digits) == 12 || len(digits) == 13) && strings.HasPrefix(digits, "55") {
		digits = digits[2:]
	}
	return len(digits) == 10 || len(digits) == 11
}

func validPostalCode(value string) bool {
	postalCode := strings.TrimSpace(value)
	if len(postalCode) == 9 && postalCode[5] == '-' {
		postalCode = postalCode[:5] + postalCode[6:]
	}
	if len(postalCode) != 8 {
		return false
	}
	for _, char := range postalCode {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// ValidateCheckout prevents reserving stock when the store cannot fulfill or
// accept payment for an order. Delivery address is only required for delivery.
func (s StoreConfig) ValidateCheckout(profile Profile) error {
	if !s.DeliveryIsActive && !s.IsPickupAvailable() {
		return ErrFulfillmentUnavailable
	}
	if !s.PaymentCashIsActive && !s.IsPixAvailable() {
		return ErrPaymentUnavailable
	}
	if !profile.HasContactDetails() {
		return ErrContactIncomplete
	}
	if !s.IsPickupAvailable() && !profile.HasDeliveryAddress() {
		return ErrDeliveryAddressIncomplete
	}
	return nil
}

func (s StoreConfig) ValidateFulfillment(profile Profile, delivery bool) error {
	if !profile.HasContactDetails() {
		return ErrContactIncomplete
	}
	if delivery {
		if !s.DeliveryIsActive {
			return ErrDeliveryUnavailable
		}
		if !profile.HasDeliveryAddress() {
			return ErrDeliveryAddressIncomplete
		}
	} else if !s.IsPickupAvailable() {
		return ErrPickupUnavailable
	}
	return nil
}
