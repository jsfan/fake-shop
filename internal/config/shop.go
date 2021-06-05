package config

import (
	"fmt"
	"github.com/jsfan/fake-shop/internal/store"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// ReadInventory reads the inventory from a YAML file
func ReadInventory(inputFile string) ([]*store.Product, error) {
	inventoryFile, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open inventory file: %w", err)
	}
	inventoryIn, err := ioutil.ReadAll(inventoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %w", err)
	}
	inventory := make([]*store.Product, 0)
	err = yaml.Unmarshal(inventoryIn, &inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inventory file: %w", err)
	}
	return inventory, nil
}

// ReadPromotions reads the promotions to be applied to purchases
func ReadPromotions(inputFile string) ([]*store.Promotion, error) {
	promotionsFile, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open promotions file: %w", err)
	}
	promotionsIn, err := ioutil.ReadAll(promotionsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read promotions file: %w", err)
	}
	promotions := make([]*store.Promotion, 0)
	err = yaml.Unmarshal(promotionsIn, &promotions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse promotions file: %w", err)
	}
	return promotions, nil
}
