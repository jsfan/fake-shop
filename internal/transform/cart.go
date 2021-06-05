package transform

import (
	"github.com/google/uuid"
	"github.com/jsfan/fake-shop/internal/graph/model"
	"github.com/jsfan/fake-shop/internal/store"
)

// RefreshCart refreshes a cart ready for delivery to the frontend
func RefreshCart(cartUUID string, cart *store.Cart) (*model.Cart, error) {
	regular, promo, errorList := cart.Get()
	outCart := &model.Cart{
		ID:             cartUUID,
		AddedItems:     nil,
		PromotionItems: nil,
		TotalPrice:     0,
		Errors:         nil,
	}
	total := 0.
	if regular != nil {
		outCart.AddedItems = make([]*model.Product, 0)
		for _, p := range regular {
			outCart.AddedItems = append(outCart.AddedItems, &model.Product{
				Sku:   p.SKU,
				Name:  p.Name,
				Price: p.Price,
				Count: &p.Count,
			})
			total += p.Price * float64(p.Count)
		}
	}
	if promo != nil {
		outCart.PromotionItems = make([]*model.Product, 0)
		for _, p := range promo {
			outCart.PromotionItems = append(outCart.PromotionItems, &model.Product{
				Sku:   p.SKU,
				Name:  p.Name,
				Price: p.Price,
				Count: &p.Count,
			})
			total += p.Price * float64(p.Count)
		}
	}
	if errorList != nil {
		outCart.Errors = make([]string, 0)
		for _, e := range errorList {
			outCart.Errors = append(outCart.Errors, e.Error())
		}
	}
	outCart.TotalPrice = total
	return outCart, nil
}

// LoadCart loads a cart with products requested from the frontend
func LoadCart(inCart model.NewCart) (*store.Cart, []error, error) {
	cartId := inCart.CartID
	cartUUID := uuid.New()
	if cartId != nil {
		var err error
		if cartUUID, err = uuid.Parse(*cartId); err != nil {
			return nil, nil, err
		}
	}
	_, outCart := store.RetrieveCart(&cartUUID)
	newItems := make([]*store.Product, 0)
	for _, p := range inCart.Products {
		newItems = append(newItems, &store.Product{
			SKU:   p.Product,
			Count: p.Count,
		})
	}
	errorList := outCart.Update(newItems)
	return outCart, errorList, nil
}

// FilterInventory filters the inventory to not contain counts
func FilterInventory() []*model.Product {
	inventory := store.GetInventory()
	filtered := make([]*model.Product, 0)
	for _, p := range inventory {
		filtered = append(filtered, &model.Product{
			Sku:   p.SKU,
			Name:  p.Name,
			Price: p.Price,
			Count: nil,
		})
	}
	return filtered
}
