package store_test

import (
	"github.com/jsfan/fake-shop/internal/store"
	"reflect"
	"testing"
)

func TestStockShop(t *testing.T) {
	faultyStock := []*store.Product{
		{
			SKU:   "A1234",
			Name:  "Carrot",
			Price: 1.1,
			Count: 10,
		},
		{
			SKU:   "A1234",
			Name:  "Stick",
			Price: 0.1,
			Count: 5,
		},
	}
	err := store.StockShop(faultyStock)
	if err == nil {
		t.Fatal("Stocking shop didn't fail for duplicate SKU.")
	}
}

func TestGetInventory(t *testing.T) {
	loadedStock := []*store.Product{
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
	expectedStock := make(map[string]*store.Product)
	for _, p := range loadedStock {
		expectedStock[p.SKU] = p
	}
	err := store.StockShop(loadedStock)
	if err != nil {
		t.Fatalf("Stocking shop failed: %+v", err)
	}
	stock := store.GetInventory()
	if !reflect.DeepEqual(stock, expectedStock) {
		t.Fatalf("Stocking shop resulted in incorrect inventory. Expected %+v, got %+v.", expectedStock, stock)
	}
}

func TestClaimInventory(t *testing.T) {
	initialStock := []*store.Product{
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
	expectedStock := make(map[string]*store.Product)
	for _, p := range initialStock {
		expectedStock[p.SKU] = p
	}
	toCart := store.Product{
		SKU:   "A1234",
		Name:  "Carrot",
		Price: 1.1,
		Count: 5,
	}
	actual, err := store.ClaimInventory(toCart)
	if err != nil {
		t.Fatalf("Failed to claim existing stock: %+v", err)
	}
	if !reflect.DeepEqual(actual, &toCart) {
		t.Errorf("Successful claim not as expected. Expected %+v, got %+v.", &toCart, actual)
	}
}
