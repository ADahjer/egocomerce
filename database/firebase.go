package database

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type Store struct {
	FireAuth  *auth.Client
	FireStore *firestore.Client
}

var Firebase *Store

func NewStore() (*Store, error) {
	path := filepath.Join("database", "egocomerce.json")
	opt := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating Auth client: %v", err)
	}

	store, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating Firebase client: %v", err)
	}

	Firebase = &Store{
		FireAuth:  client,
		FireStore: store,
	}

	return Firebase, nil
}
