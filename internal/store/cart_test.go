package store_test

import (
	"github.com/jsfan/fake-shop/internal/store"
	"reflect"
	"testing"
)

func setupShop() (*store.Cart, error) {
	stock := []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 10,
		},
		{
			SKU:   "B1234",
			Name:  "Stick",
			Price: 0.1,
			Count: 5,
		},
	}
	if err := store.StockShop(stock); err != nil {
		return nil, err
	}
	return &store.Cart{}, nil
}

func TestCart_Add(t *testing.T) {
	c, err := setupShop()
	if err != nil {
		t.Fatalf("Test setup failed: %+v", err)
	}
	err = c.Add(&store.Product{
		SKU:   "A1234",
		Name:  "Carrot",
		Price: 1.1,
		Count: 5,
	})
	if err != nil {
		t.Errorf("Unexpected error when adding to cart: %+v", err)
	}
	err = c.Add(&store.Product{
		SKU:   "A1234",
		Name:  "Carrot",
		Price: 1.1,
		Count: 5,
	})
	if err != nil {
		t.Errorf("Unexpected error when adding to cart: %+v", err)
	}
	err = c.Add(&store.Product{
		SKU:   "A1234",
		Name:  "Carrot",
		Price: 1.1,
		Count: 1,
	})
	if err == nil || err.Error() != "not enough stock" {
		t.Errorf("Unexpected error when adding to cart: %+v", err)
	}
}

func TestCart_Update(t *testing.T) {
	c, err := setupShop()
	if err != nil {
		t.Fatalf("Test setup failed: %+v", err)
	}
	newCart := []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 10,
		},
		{
			SKU:   "B1234",
			Name:  "Carrot",
			Price: 0.1,
			Count: 5,
		},
	}
	errors := c.Update(newCart)
	if errors != nil {
		t.Fatalf("Updating cart failed unexpectedly: %+v", errors)
	}

	newCart = []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 5,
		},
		{
			SKU:   "B1234",
			Name:  "Carrot",
			Price: 0.1,
			Count: 3,
		},
	}
	errors = c.Update(newCart)
	if errors != nil {
		t.Fatalf("Updating cart failed unexpectedly: %+v", errors)
	}
	newCart = []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 5,
		},
		{
			SKU:   "B1234",
			Name:  "Carrot",
			Price: 0.1,
			Count: 10,
		},
	}
	errors = c.Update(newCart)
	if len(errors) == 0 {
		t.Fatal("Updating cart with claim for non-existent stock threw no error.")
	}
	if len(errors) > 1 {
		t.Fatalf("Updating cart with claim for non-existent stock threw unexpected additional errors: %+v", errors)
	}
	if errors[0].Error() != "not enough stock" {
		t.Fatalf(`Did not get expected out of stock error. Expected "not enough stock", got "%+v".`, errors[0])
	}
}

func TestCart_Get(t *testing.T) {
	c, err := setupShop()
	if err != nil {
		t.Fatalf("Test setup failed: %+v", err)
	}
	promos := []*store.Promotion{
		{
			Name:     "A freebie",
			SKU:      "FREEBIE",
			Category: "freebie",
			Requires: store.Requirement{
				SKU:   "A1234",
				Count: 1,
			},
			Rule: store.RuleDetail{
				SKU:      "B1234",
				Count:    1,
				Discount: 0,
			},
		},
	}
	newCart := []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 1,
		},
		{
			SKU:   "B1234",
			Name:  "Carrot",
			Price: 0.1,
			Count: 3,
		},
	}
	expectedCart := make(map[string]*store.Product, 0)
	for _, p := range newCart {
		expectedCart[p.SKU] = p
	}
	errors := c.Update(newCart)
	if errors != nil {
		t.Fatalf("Initialising cart failed unexpectedly: %+v", errors)
	}
	cart, promo, errors := c.Get()
	if !reflect.DeepEqual(cart, expectedCart) {
		t.Errorf("Retrieved cart does not contain expected items. Expected %+v, got +%v.", expectedCart, cart)
	}
	if len(promo) != 0 {
		t.Errorf("Got unexpected promotion items: %+v", promo)
	}
	if errors != nil {
		t.Errorf("Got unexpected error when retrieving cart: %+v", errors)
	}

	expectedPromo := map[string]*store.Product{
		"FREEBIE": {
			SKU:   "FREEBIE",
			Name:  "A freebie",
			Price: 0,
			Count: 1,
		},
	}
	store.RegisterPromotions(promos)
	cart, promo, errors = c.Get()
	if !reflect.DeepEqual(cart, expectedCart) {
		t.Errorf("Retrieved cart does not contain expected items. Expected %+v, got +%v.", expectedCart, cart)
	}
	if !reflect.DeepEqual(promo, expectedPromo) {
		t.Errorf("Got unexpected promotion items. Expected %+v, got %+v", expectedPromo, promo)
	}
	if errors != nil {
		t.Errorf("Got unexpected error when retrieving cart: %+v", errors)
	}
	store.RegisterPromotions(promos)
}
