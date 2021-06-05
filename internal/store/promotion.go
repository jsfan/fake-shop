package store

import "fmt"

var promotions []*Promotion

type Requirement struct {
	SKU   string
	Count int
}

type RuleDetail struct {
	SKU      string
	Count    int
	Discount float64
}

type Promotion struct {
	Name     string
	SKU      string
	Category string
	Requires Requirement `yaml:"requires"`
	Rule     RuleDetail  `yaml:"rule"`
}

// RegisterPromotions takes a list of promotions and registers them for use
func RegisterPromotions(promos []*Promotion) {
	promotions = promos
}

// Apply applies a promotion to a product
func (p *Promotion) Apply(product *Product) (claimsItem *Product, promoItem *Product, err error) {
	if product.SKU == p.Requires.SKU {
		switch p.Category {
		case "freebie":
			freebieCount := product.Count / p.Requires.Count
			return &Product{
					SKU:   p.Rule.SKU,
					Name:  "", // filled later by cart using SKU
					Price: 0,
					Count: freebieCount,
				},
				&Product{
					SKU:   p.SKU,
					Name:  p.Name,
					Price: 0,
					Count: freebieCount,
				},
				nil
		case "n4m":
			freeItemCount := product.Count/p.Requires.Count*p.Requires.Count - p.Rule.Count
			return nil, &Product{
				SKU:   p.SKU,
				Name:  p.Name,
				Price: -product.Price,
				Count: freeItemCount,
			}, nil
		case "discount":
			if product.Count >= p.Requires.Count {
				return nil, &Product{
					SKU:   p.SKU,
					Name:  p.Name,
					Price: -p.Rule.Discount * product.Price,
					Count: product.Count,
				}, nil
			}
		default:
			return product, nil, fmt.Errorf(`unknown promotion "%s"`, p.Category)
		}
	}
	return nil, nil, nil
}
