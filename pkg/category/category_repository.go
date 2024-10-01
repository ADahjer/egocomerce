package category

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/types"
	"google.golang.org/api/iterator"
)

var s *database.Store

const collectionName = "categories"

func InitRepo() {
	s = database.Firebase
}

func CreateCategory(ctx context.Context, name string) (string, error) {
	docRef, _, err := s.FireStore.Collection(collectionName).Add(ctx, &CreateCategoryModel{Name: name})
	if err != nil {
		return "", err
	}

	return docRef.ID, nil

}

func GetAllCategories(ctx context.Context) (interface{}, error) {
	iter := s.FireStore.Collection(collectionName).Documents(ctx)

	defer iter.Stop()

	var categories []types.Map

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		ctg := types.Map{
			"id":   doc.Ref.ID,
			"name": doc.Data()["Name"],
		}

		categories = append(categories, ctg)
	}

	return categories, nil
}

func GetCategoryById(ctx context.Context, id string) (interface{}, error) {
	docRef, err := s.FireStore.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	ctg := types.Map{
		"id":   docRef.Ref.ID,
		"name": docRef.Data()["Name"],
	}

	return ctg, nil

}

func DeleteCategory(ctx context.Context, id string) (bool, error) {
	_, err := s.FireStore.Collection(collectionName).Doc(id).Delete(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateCategory(ctx context.Context, id, name string) (interface{}, error) {

	docRef := s.FireStore.Collection(collectionName).Doc(id)

	_, err := docRef.Set(ctx, types.Map{"Name": name}, firestore.MergeAll)
	if err != nil {
		return nil, err
	}

	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	ctg := types.Map{
		"id":   doc.Ref.ID,
		"name": name,
	}

	return ctg, nil
}
