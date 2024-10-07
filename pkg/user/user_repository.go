package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/types"
)

var s *database.Store

func InitRepo() {
	s = database.Firebase
}

func CreateUser(ctx context.Context, userName, email, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		DisplayName(userName).
		Email(email).
		Password(password)

	userRecord, err := s.FireAuth.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}

func LoginWithEmailAndPassword(email, password string) (types.Map, error) {
	apiKey := os.Getenv("FIREBASE_API_KEY")
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	payload := types.Map{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login %v", result["error"])
	}

	token, ok := result["idToken"].(string)
	id := result["localId"].(string)
	log.Println(id)

	userInfo, err := GetUSerInfo(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("unable to get token")
	}

	tokenWithClaims := types.Map{
		"token":  token,
		"claims": userInfo.CustomClaims,
	}

	return tokenWithClaims, nil

}

func VerifyToken(ctx context.Context, token string) (*auth.Token, error) {
	decodedToken, err := s.FireAuth.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return decodedToken, nil
}

func GetUSerInfo(ctx context.Context, uid string) (*auth.UserRecord, error) {
	userRecord, err := s.FireAuth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}
