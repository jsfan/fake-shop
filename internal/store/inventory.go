package store

import (
	"fmt"
)

var inventory map[string]*Product

type Product struct {
	SKU   string
	Name  string
	Price float64
	Count int `yaml:"stock"`
}

// StockShop takes an inventory and stocks the shop with it
func StockShop(stock []*Product) error {
	inventory = make(map[string]*Product, 0)
	for _, item := range stock {
		if _, ok := inventory[item.SKU]; ok {
			return fmt.Errorf(`found duplicate SKU "%s"`, item.SKU)
		}
		inventory[item.SKU] = item
	}
	return nil
}

// ClaimInventory claims stock from the inventory to add to a cart
func ClaimInventory(product Product) (*Product, error) {
	successfulClaim := product
	if _, ok := inventory[product.SKU]; !ok {
		successfulClaim.Count = 0
		return nil, fmt.Errorf(`SKU "%s" does not exist`, product.SKU)
	}
	invProd := inventory[product.SKU]
	invProd.Count -= product.Count
	if invProd.Count < 0 {
		successfulClaim.Count += inventory[product.SKU].Count
		inventory[product.SKU].Count = 0
		return &successfulClaim, fmt.Errorf(`not enough stock`)
	}
	successfulClaim.Name = invProd.Name
	successfulClaim.Price = invProd.Price
	return &successfulClaim, nil
}

// GetInventory returns the current inventory
func GetInventory() map[string]*Product {
	return inventory
}
