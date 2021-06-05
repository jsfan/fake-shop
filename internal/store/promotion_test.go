package store_test

import (
	"github.com/jsfan/fake-shop/internal/store"
	"reflect"
	"testing"
)

func TestPromotion_ApplyFreebie(t *testing.T) {
	promo := &store.Promotion{
		Name:     "A freebie",
		SKU:      "FREEBIE",
		Category: "freebie",
		Requires: store.Requirement{
			SKU:   "ABC123",
			Count: 1,
		},
		Rule: store.RuleDetail{
			SKU:      "DEF567",
			Count:    1,
			Discount: 0,
		},
	}
	prod := &store.Product{
		SKU:   "12345678",
		Name:  "Test",
		Price: 11.5,
		Count: 2,
	}
	claims, promoItems, err := promo.Apply(prod)
	if claims != nil {
		t.Errorf("Got a claim for a non-matching promotion.")
	}
	if promoItems != nil {
		t.Errorf("Got promotion items for a non-matching promotion.")
	}
	if err != nil {
		t.Errorf("Got error for non-matching promotion: %+v", err)
	}
	prod = &store.Product{
		SKU:   "ABC123",
		Name:  "Test",
		Price: 11.5,
		Count: 2,
	}
	expectedClaims := &store.Product{
		SKU:   "DEF567",
		Name:  "",
		Price: 0,
		Count: 2,
	}
	expectedPromoItems := &store.Product{
		SKU:   "FREEBIE",
		Name:  "A freebie",
		Price: 0,
		Count: 2,
	}
	claims, promoItems, err = promo.Apply(prod)
	if !reflect.DeepEqual(claims, expectedClaims) {
		t.Errorf("Incorrect claim. Expected %+v, got %+v.", expectedClaims, claims)
	}
	if !reflect.DeepEqual(promoItems, expectedPromoItems) {
		t.Errorf("Incorrect promotion items. Expected %+v, got %+v.", expectedPromoItems, promoItems)
	}
	if err != nil {
		t.Errorf("Got unexpected error for promotion: %+v", err)
	}
}

func TestPromotion_ApplyDiscount(t *testing.T) {
	promo := &store.Promotion{
		Name:     "A discount",
		SKU:      "DISCOUNT",
		Category: "discount",
		Requires: store.Requirement{
			SKU:   "ABC123",
			Count: 10,
		},
		Rule: store.RuleDetail{
			Discount: .3,
		},
	}
	prod := &store.Product{
		SKU:   "12345678",
		Name:  "Test",
		Price: 11.5,
		Count: 2,
	}
	claims, promoItems, err := promo.Apply(prod)
	if claims != nil {
		t.Errorf("Got a claim for a non-matching promotion.")
	}
	if promoItems != nil {
		t.Errorf("Got promotion items for a non-matching promotion.")
	}
	if err != nil {
		t.Errorf("Got error for non-matching promotion: %+v", err)
	}
	prod = &store.Product{
		SKU:   "ABC123",
		Name:  "Test",
		Price: 15.,
		Count: 2,
	}
	claims, promoItems, err = promo.Apply(prod)
	if claims != nil {
		t.Errorf("Got a claim for promotion with unmet requirements.")
	}
	if promoItems != nil {
		t.Errorf("Got promotion for promotion with unmet requirements.")
	}

	prod.Count = 10
	expectedPromoItems := &store.Product{
		SKU:   "DISCOUNT",
		Name:  "A discount",
		Price: -4.5,
		Count: 10,
	}
	claims, promoItems, err = promo.Apply(prod)
	if err != nil {
		t.Errorf("Got error for for promotion with unmet requirements: %+v", err)
	}
	if claims != nil {
		t.Errorf("Got a claim for a discount promotion: %+v", claims)
	}
	if !reflect.DeepEqual(promoItems, expectedPromoItems) {
		t.Errorf("Incorrect promotion items. Expected %+v, got %+v.", expectedPromoItems, promoItems)
	}
	if err != nil {
		t.Errorf("Got unexpected error for promotion: %+v", err)
	}
}
func TestPromotion_Apply(t *testing.T) {
	promo := &store.Promotion{
		Name:     "A 2 for 1",
		SKU:      "2FOR1",
		Category: "n4m",
		Requires: store.Requirement{
			SKU:   "ABC123",
			Count: 2,
		},
		Rule: store.RuleDetail{
			Count: 1,
		},
	}
	prod := &store.Product{
		SKU:   "12345678",
		Name:  "Test",
		Price: 11.5,
		Count: 2,
	}
	claims, promoItems, err := promo.Apply(prod)
	if claims != nil {
		t.Errorf("Got a claim for a non-matching promotion.")
	}
	if promoItems != nil {
		t.Errorf("Got promotion items for a non-matching promotion.")
	}
	if err != nil {
		t.Errorf("Got error for non-matching promotion: %+v", err)
	}
	prod = &store.Product{
		SKU:   "ABC123",
		Name:  "Test",
		Price: 11.5,
		Count: 2,
	}
	expectedPromoItems := &store.Product{
		SKU:   "2FOR1",
		Name:  "A 2 for 1",
		Price: -11.5,
		Count: 1,
	}
	claims, promoItems, err = promo.Apply(prod)
	if claims != nil {
		t.Errorf("Got a claim for an N for M promotion: %+v", claims)
	}
	if !reflect.DeepEqual(promoItems, expectedPromoItems) {
		t.Errorf("Incorrect promotion items. Expected %+v, got %+v.", expectedPromoItems, promoItems)
	}
	if err != nil {
		t.Errorf("Got unexpected error for promotion: %+v", err)
	}
}
