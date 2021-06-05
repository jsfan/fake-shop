package config_test

import (
	"github.com/jsfan/fake-shop/internal/config"
	"github.com/jsfan/fake-shop/internal/store"
	"reflect"
	"testing"
)

func TestReadInventory(t *testing.T) {
	expectedStock := []store.Product{
		{
			SKU:   "1234",
			Name:  "Some Item",
			Price: 9.99,
			Count: 12,
		},
		{
			SKU:   "ABC123",
			Name:  "Another Item",
			Price: 123.45,
			Count: 22,
		},
	}
	_, err := config.ReadInventory("missing.yaml")
	if err == nil {
		t.Error("Loading missing file did not throw an error.")
	} else if err.Error()[:30] != "failed to open inventory file:" {
		t.Errorf("Got unexpected error when loading incorrectly formatted YAML file: %+v", err)
	}
	_, err = config.ReadInventory("../../test/data/bad_stock.yaml")
	if err == nil {
		t.Error("Loading an incorrectly formatted YAML file did not throw an error.")
	} else if err.Error()[:31] != "failed to parse inventory file:" {
		t.Errorf("Got unexpected error when loading incorrectly formatted YAML file: %+v", err)
	}
	stock, err := config.ReadInventory("../../test/data/good_stock.yaml")
	if err != nil {
		t.Errorf("Got an unexpected error when loading good stock file: %+v", err)
	}
	stockCopy := make([]store.Product, 0)
	for _, p := range stock {
		stockCopy = append(stockCopy, *p)
	}
	if !reflect.DeepEqual(stockCopy, expectedStock) {
		t.Errorf("Loaded stock is not as expected. Expected %+v, got %+v.", expectedStock, stockCopy)
	}
}

func TestReadPromotions(t *testing.T) {
	expectedPromotions := []store.Promotion{
		{
			SKU:      "FREEBIE4U",
			Name:     "A freebie",
			Category: "freebie",
			Requires: store.Requirement{
				SKU:   "ABC123",
				Count: 1,
			},
			Rule: store.RuleDetail{
				SKU:      "12345",
				Count:    1,
				Discount: 0,
			},
		},
		{
			SKU:      "100OFF",
			Name:     "A discount",
			Category: "discount",
			Requires: store.Requirement{
				SKU:   "98765",
				Count: 10,
			},
			Rule: store.RuleDetail{
				SKU:      "",
				Count:    0,
				Discount: 1.,
			},
		},
	}
	_, err := config.ReadPromotions("missing.yaml")
	if err == nil {
		t.Error("Loading missing file did not throw an error.")
	} else if err.Error()[:31] != "failed to open promotions file:" {
		t.Errorf("Got unexpected error when loading incorrectly formatted YAML file: %+v", err)
	}
	_, err = config.ReadPromotions("../../test/data/bad_stock.yaml")
	if err == nil {
		t.Error("Loading an incorrectly formatted YAML file did not throw an error.")
	} else if err.Error()[:32] != "failed to parse promotions file:" {
		t.Errorf("Got unexpected error when loading incorrectly formatted YAML file: %+v", err)
	}
	promotions, err := config.ReadPromotions("../../test/data/good_promotions.yaml")
	if err != nil {
		t.Errorf("Got an unexpected error when loading good promotions file: %+v", err)
	}
	promoCopy := make([]store.Promotion, 0)
	for _, p := range promotions {
		promoCopy = append(promoCopy, *p)
	}
	if !reflect.DeepEqual(promoCopy, expectedPromotions) {
		t.Errorf("Loaded promotions are not as expected. Expected %+v, got %+v.", expectedPromotions, promoCopy)
	}
}
