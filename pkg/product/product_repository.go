package product

import (
	"context"

	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/pkg/category"
	"google.golang.org/api/iterator"
)

var s *database.Store

const collectionName = "products"

func InitRepo() {
	s = database.Firebase
}

func CreateProduct(ctx context.Context, product CreateProductModel) (string, error) {

	for _, c := range product.Categories {
		_, err := category.GetCategoryById(ctx, c)
		if err != nil {
			return "", err
		}
	}

	docRef, _, err := s.FireStore.Collection(collectionName).Add(ctx, &CreateProductModel{
		Name:       product.Name,
		Price:      product.Price,
		Categories: product.Categories,
	})
	if err != nil {
		return "", err
	}

	return docRef.ID, nil
}

func GetAllProducts(ctx context.Context) ([]ProductModel, error) {
	iter := s.FireStore.Collection(collectionName).Documents(ctx)

	defer iter.Stop()

	var products []ProductModel

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		prod := ProductModel{
			Id:         doc.Ref.ID,
			Name:       doc.Data()["Name"].(string),
			Price:      doc.Data()["Price"].(float64),
			Categories: doc.Data()["Categories"],
		}

		products = append(products, prod)
	}

	return products, nil

}

func GetProductById(ctx context.Context, id string) (*ProductModel, error) {
	doc, err := s.FireStore.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	prod := &ProductModel{
		Id:         doc.Ref.ID,
		Name:       doc.Data()["Name"].(string),
		Price:      doc.Data()["Price"].(float64),
		Categories: doc.Data()["Categories"],
	}

	return prod, nil
}
