package models

import (
	"slices"
	"testing"
)

func TestAvailableTransitionsFollowOrderRules(t *testing.T) {
	for _, delivery := range []bool{false, true} {
		for _, current := range []ShoppingCartStatus{ActiveStatus, AwaitingPaymentStatus, PaymentApprovedStatus, ProcessingStatus, DeliveredStatusAwaiting, DeliveredStatusSent, DeliveredStatusDelivered, CancelledStatus} {
			order := Order{Status: current, IsDelivery: delivery}
			for _, next := range order.AvailableTransitions() {
				if next == current || !order.CanTransitionTo(next) {
					t.Fatalf("offered invalid transition %q -> %q", current, next)
				}
			}
		}
	}
}

func TestPickupCanBeCollectedButCannotBeShipped(t *testing.T) {
	pickup := Order{Status: DeliveredStatusAwaiting}
	if pickup.CanTransitionTo(DeliveredStatusSent) || slices.Contains(pickup.AvailableTransitions(), DeliveredStatusSent) {
		t.Fatal("pickup must not offer shipment")
	}
	if !slices.Contains(pickup.AvailableTransitions(), DeliveredStatusDelivered) {
		t.Fatal("pickup must offer collection")
	}
	delivery := Order{Status: DeliveredStatusAwaiting, IsDelivery: true}
	if !slices.Contains(delivery.AvailableTransitions(), DeliveredStatusSent) {
		t.Fatal("delivery must offer shipment")
	}
	if slices.Contains(delivery.AvailableTransitions(), DeliveredStatusDelivered) {
		t.Fatal("delivery must be shipped before delivery is confirmed")
	}
}
