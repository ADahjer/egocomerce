package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"github.com/ADahjer/egocomerce/types"

	"google.golang.org/api/option"
)

type Store struct {
	FireAuth    *auth.Client
	FireStore   *firestore.Client
	FireStorage *storage.Client
}

var Firebase *Store

func NewStore() (*Store, error) {
	firebaseKey := os.Getenv("FIREBASE_SERVICE_KEY")

	// Get the downloaded credential.json from json from an env variable and parse it into a json
	var credentialMap types.Map
	if err := json.Unmarshal([]byte(firebaseKey), &credentialMap); err != nil {
		return nil, err
	}

	credsJson, err := json.Marshal(credentialMap)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credsJson)

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

	storage, err := app.Storage(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error creating fibase storage: %v", err)
	}

	Firebase = &Store{
		FireAuth:    client,
		FireStore:   store,
		FireStorage: storage,
	}

	return Firebase, nil
}
