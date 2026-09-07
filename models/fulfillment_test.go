package models

import (
	"errors"
	"testing"
)

func TestPickupRequiresUsableStoreAddress(t *testing.T) {
	for _, tc := range []struct {
		address, city, state string
		available            bool
	}{
		{"", "", "", false}, {"   ", "São Paulo", "SP", false},
		{"Rua 1", "", "SP", false}, {"Rua 1", "São Paulo", "XX", false},
		{"Rua 1", "São Paulo", "São Paulo", false}, {" Rua 1 ", " São Paulo ", " sp ", true},
	} {
		store := StoreConfig{PhysicalStoreAddress: tc.address, PhysicalStoreCity: tc.city, PhysicalStoreState: tc.state}
		if got := store.IsPickupAvailable(); got != tc.available {
			t.Fatalf("%+v: available=%v", tc, got)
		}
	}
}

func TestContactAndDeliveryAddressAreSeparateRequirements(t *testing.T) {
	profile := Profile{FirstName: " Cliente ", LastName: " Teste ", PhoneNumber: "+55 (11) 99999-9999"}
	if !profile.HasContactDetails() || profile.HasDeliveryAddress() || profile.IsProfileComplete() {
		t.Fatal("pickup contact incorrectly requires address")
	}
	for _, field := range []string{"name", "surname", "phone"} {
		invalid := profile
		switch field {
		case "name":
			invalid.FirstName = " "
		case "surname":
			invalid.LastName = " "
		case "phone":
			invalid.PhoneNumber = "telefone"
		}
		if invalid.HasContactDetails() {
			t.Fatalf("accepted invalid %s", field)
		}
	}
	profile.Address, profile.City, profile.State, profile.PostalCode = "Rua 1", "São Paulo", " sp ", "01001-000"
	if !profile.IsProfileComplete() {
		t.Fatal("valid address rejected")
	}
	for _, postalCode := range []string{"", "123", "abcdefgh", "+01001000", "01001 000"} {
		profile.PostalCode = postalCode
		if profile.HasDeliveryAddress() {
			t.Fatalf("accepted invalid CEP %q", postalCode)
		}
	}
}

func TestFulfillmentNeverSilentlyConvertsTheCustomersChoice(t *testing.T) {
	profile := Profile{FirstName: "Cliente", LastName: "Teste", PhoneNumber: "11999999999"}
	store := StoreConfig{PaymentCashIsActive: true, PhysicalStoreAddress: "Rua 1", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP"}
	if err := store.ValidateFulfillment(profile, false); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateFulfillment(profile, true); !errors.Is(err, ErrDeliveryUnavailable) {
		t.Fatalf("unexpected error: %v", err)
	}
	store.DeliveryIsActive = true
	if err := store.ValidateFulfillment(profile, true); !errors.Is(err, ErrDeliveryAddressIncomplete) {
		t.Fatalf("unexpected error: %v", err)
	}
	store.PhysicalStoreAddress = ""
	if err := store.ValidateFulfillment(profile, false); !errors.Is(err, ErrPickupUnavailable) {
		t.Fatalf("unexpected error: %v", err)
	}
}
