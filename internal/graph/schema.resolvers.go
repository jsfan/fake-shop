package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jsfan/fake-shop/internal/graph/generated"
	"github.com/jsfan/fake-shop/internal/graph/model"
	"github.com/jsfan/fake-shop/internal/store"
	"github.com/jsfan/fake-shop/internal/transform"
)

func (r *mutationResolver) AddProduct(ctx context.Context, input model.AdditionalItem) (*model.Cart, error) {
	cartUUID := uuid.New()
	var err error
	if input.CartID != nil {
		cartUUID, err = uuid.Parse(*input.CartID)
		if err != nil {
			return nil, errors.New("invalid Cart ID")
		}
	}
	_, cart := store.RetrieveCart(&cartUUID)
	newProduct := &store.Product{
		SKU:   input.Item.Product,
		Count: input.Item.Count,
	}
	err = cart.Add(newProduct)
	if err != nil {
		return nil, err
	}
	return transform.RefreshCart(cartUUID.String(), cart)
}

func (r *mutationResolver) UpdateCart(ctx context.Context, input model.NewCart) (*model.Cart, error) {
	cartUUID := uuid.New()
	var err error
	if input.CartID != nil {
		cartUUID, err = uuid.Parse(*input.CartID)
		if err != nil {
			return nil, errors.New("invalid Cart ID")
		}
	}
	_, cart := store.RetrieveCart(&cartUUID)
	cart, errorList, err := transform.LoadCart(input)
	if err != nil {
		return nil, err
	}
	outCart, err := transform.RefreshCart(cartUUID.String(), cart)
	if err != nil {
		return nil, err
	}
	errorStrings := make([]string, 0)
	if errorList != nil {
		for _, e := range errorList {
			errorStrings = append(errorStrings, e.Error())
		}
	}
	if outCart.Errors != nil {
		outCart.Errors = append(errorStrings, outCart.Errors...)
	}
	return outCart, nil
}

func (r *queryResolver) Cart(ctx context.Context, input *string) (*model.Cart, error) {
	cartUUID := uuid.New()
	var err error
	if input != nil {
		cartUUID, err = uuid.Parse(*input)
		if err != nil {
			return nil, errors.New("invalid Cart ID")
		}
	}
	_, cart := store.RetrieveCart(&cartUUID)
	return transform.RefreshCart(cartUUID.String(), cart)
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	return transform.FilterInventory(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
