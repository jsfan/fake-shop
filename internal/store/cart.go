package store

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	contents   map[string]*Product
	promoCache map[string]*Product
	expires    time.Time
}

var carts map[uuid.UUID]*Cart

func InitShop() {
	carts = make(map[uuid.UUID]*Cart, 0)
}

// Add adds a product to a cart with an item count
func (c *Cart) Add(product *Product) error {
	if c.contents == nil {
		c.contents = make(map[string]*Product)
	}
	claims, err := ClaimInventory(*product)
	if err != nil {
		return err
	}
	if inCart, ok := c.contents[product.SKU]; ok { // add new item
		inCart.Count += claims.Count
		c.contents[product.SKU] = inCart
	} else { // add existing
		c.contents[product.SKU] = claims
	}
	return err
}

// Update replaces the cart contents with those submitted
func (c *Cart) Update(products []*Product) []error {
	errors := make([]error, 0)
	if c.contents == nil {
		c.contents = make(map[string]*Product)
	}
	for _, p := range products {
		var prev *Product
		var ok bool
		if prev, ok = c.contents[p.SKU]; ok {
			p = &Product{
				SKU:   p.SKU,
				Name:  p.Name,
				Price: p.Price,
				Count: p.Count - prev.Count,
			}
		} else {
			prev = &Product{}
			*prev = *p
			prev.Count = 0
		}
		actual, err := ClaimInventory(*p)
		if err != nil {
			errors = append(errors, err)
		}
		c.contents[p.SKU] = prev
		c.contents[p.SKU].Name = actual.Name
		c.contents[p.SKU].Price = actual.Price
		c.contents[p.SKU].Count += actual.Count
	}
	if len(errors) == 0 {
		errors = nil
	}
	return errors
}

// Get retrieves the cart with promotions applied
func (c *Cart) Get() (cartItems, promoItems map[string]*Product, errors []error) {
	errors = make([]error, 0)
	if c.promoCache == nil {
		c.promoCache = make(map[string]*Product)
	}
	promoItems = make(map[string]*Product, 0)
	for _, p := range c.contents {
		for _, promo := range promotions {
			inventoryClaim, extra, err := promo.Apply(p)
			if err != nil {
				errors = append(errors, fmt.Errorf(`internal error: %w`, err))
			}
			if inventoryClaim != nil { // this promotion claims extra stock
				// check if we already have claimed items for promotion
				var cached *Product
				var ok bool
				if cached, ok = c.promoCache[inventoryClaim.SKU]; ok {
					inventoryClaim.Count = inventoryClaim.Count - cached.Count
				} else {
					cached = &Product{
						Count: 0,
					}
				}
				actual, err := ClaimInventory(*inventoryClaim)
				if err != nil {
					errors = append(errors, fmt.Errorf(`promotion could not be applied: %s`, err))
				}
				actual.Count += cached.Count
				c.promoCache[actual.SKU] = actual
				extra.Count = actual.Count
			}
			if extra != nil {
				promoItems[extra.SKU] = extra
			}
		}
	}
	if len(errors) == 0 { // no errors, so we return a null pointer
		errors = nil
	}
	return c.contents, promoItems, errors
}

// RetrieveCart retrieves a cart from memory or creates a new one
func RetrieveCart(cartId *uuid.UUID) (*uuid.UUID, *Cart) {
	var cart *Cart
	var exists bool
	if cartId != nil {
		cart, exists = carts[*cartId]
	} else {
		exists = false
		newId := uuid.New()
		cartId = &newId
	}
	if !exists {
		cart = &Cart{
			contents:   nil,
			promoCache: nil,
			expires:    time.Now().Add(600 * time.Second),
		}
		carts[*cartId] = cart
	}
	return cartId, cart
}
