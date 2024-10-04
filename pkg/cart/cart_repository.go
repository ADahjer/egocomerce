package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/ADahjer/egocomerce/database"
)

var s *database.Store

const collectionName = "carts"

func InitRepo() {
	s = database.Firebase
}

func CreateNewCart(ctx context.Context, userId string) (*CartModel, string, error) {
	newCart := &CartModel{
		UserID:    userId,
		Items:     []CartItemModel{},
		Status:    "active",
		CreatedAt: time.Now(),
	}

	ref, _, err := s.FireStore.Collection(collectionName).Add(ctx, newCart)
	if err != nil {
		return nil, "", nil
	}

	return newCart, ref.ID, nil
}

func GetCartById(ctx context.Context, cartId string) (*CartModel, error) {
	docRef, err := s.FireStore.Collection(collectionName).Doc(cartId).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting the cart: %+v", err)
	}

	var cart CartModel
	err = docRef.DataTo(&cart)
	if err != nil {
		return nil, fmt.Errorf("error parsing cart data: %+v", err)
	}

	return &cart, nil
}

// Get the active cart for the current user and return it with the id as string
func getActiveCart(ctx context.Context, userId string) (*CartModel, string, error) {
	query := s.FireStore.Collection(collectionName).Where("UserID", "==", userId).Where("Status", "==", "active").Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, "", fmt.Errorf("no active cart found for the current user")
	}

	if len(docs) == 0 {
		return nil, "", nil
	}

	var cart CartModel
	ref := docs[0]
	err = ref.DataTo(&cart)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing cart data: %+v", err)
	}

	return &cart, ref.Ref.ID, nil
}

func AddItemToCart(ctx context.Context, userId string, newItem CartItemModel) error {
	cart, cartId, err := getActiveCart(ctx, userId)
	if err != nil {
		return err
	}

	// TODO: check if the product that will be added exists

	productFound := false
	for i, item := range cart.Items {
		if item.ProductID == newItem.ProductID {
			cart.Items[i].Quantity += newItem.Quantity
			cart.Items[i].Price += newItem.Price
			productFound = true

			// TODO: Update the func to also reduce products, check if the quantity its <= 0
			break
		}
	}

	if !productFound {
		cart.Items = append(cart.Items, newItem)
	}

	cartDoc := s.FireStore.Collection(collectionName).Doc(cartId)
	_, err = cartDoc.Set(ctx, cart)
	if err != nil {
		return fmt.Errorf("could not update the cart: %+v", err)
	}

	return nil
}
