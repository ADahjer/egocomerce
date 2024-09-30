package database

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type Store struct {
	FireApp  *firebase.App
	FireAuth *auth.Client
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

	Firebase = &Store{
		FireApp:  app,
		FireAuth: client,
	}

	return Firebase, nil
}
