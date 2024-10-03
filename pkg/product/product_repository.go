package product

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/pkg/category"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

var s *database.Store

const collectionName = "products"

func InitRepo() {
	s = database.Firebase
}

func UploadProductImage(ctx context.Context, image multipart.File) (string, error, *storage.ObjectHandle) {
	imageName := fmt.Sprintf("products/%s", uuid.New().String())
	bucketName := os.Getenv("STORAGE_BUCKET")
	bucket, err := s.FireStorage.Bucket(bucketName)
	if err != nil {
		return "", err, nil
	}

	wc := bucket.Object(imageName).NewWriter(ctx)
	if _, err := io.Copy(wc, image); err != nil {
		return "", errors.New("error 3"), nil
	}

	if err := wc.Close(); err != nil {
		return "", errors.New("error 4"), nil
	}

	o := bucket.Object(imageName)
	if err := o.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", errors.New("error 5"), nil
	}

	imageUrl := fmt.Sprintf("http://storage.googleapis.com/%s/%s", bucketName, imageName)

	return imageUrl, nil, o
}

func CreateProduct(ctx context.Context, product CreateProductModel, image multipart.File) (string, error) {

	for _, c := range product.Categories {
		_, err := category.GetCategoryById(ctx, c)
		if err != nil {
			return "", err
		}
	}

	imageName, err, o := UploadProductImage(ctx, image)
	if err != nil {
		return "", err
	}

	docRef, _, err := s.FireStore.Collection(collectionName).Add(ctx, &ProductModel{
		Image:      imageName,
		Name:       product.Name,
		Price:      product.Price,
		Categories: product.Categories,
	})
	if err != nil {
		if delErr := o.Delete(ctx); delErr != nil {
			return "", err
		}
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
