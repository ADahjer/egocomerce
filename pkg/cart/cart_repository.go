package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/pkg/product"
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

func AddItemToCart(ctx context.Context, userId string, newItem NewCartItemModel) error {
	cart, cartId, err := getActiveCart(ctx, userId)
	if err != nil {
		return err
	}

	// check if the product that will be added exists
	prod, err := product.GetProductById(ctx, newItem.ProductID)
	if err != nil {
		return err
	}

	newPrice := prod.Price * float64(newItem.Quantity)

	productFound := false
	for i, item := range cart.Items {
		if item.ProductID == newItem.ProductID {
			productFound = true
			newQuantity := cart.Items[i].Quantity + newItem.Quantity

			// if there will be no more of that item left, just remove it
			if newQuantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)

			} else {
				cart.Items[i].Quantity += newItem.Quantity
				cart.Items[i].Price += newPrice
			}
			// TODO: Update the func to also reduce products, check if the quantity its <= 0
			break
		}
	}

	if !productFound {
		itemToInsert := &CartItemModel{
			ProductID: newItem.ProductID,
			Quantity:  newItem.Quantity,
			Price:     newPrice,
		}
		cart.Items = append(cart.Items, *itemToInsert)
	}

	cartDoc := s.FireStore.Collection(collectionName).Doc(cartId)
	_, err = cartDoc.Set(ctx, cart)
	if err != nil {
		return fmt.Errorf("could not update the cart: %+v", err)
	}

	return nil
}
