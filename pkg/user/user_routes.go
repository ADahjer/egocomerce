package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ADahjer/egocomerce/types"
	"github.com/ADahjer/egocomerce/utils"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Group) {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)
	router.GET("/profile", handleGetProfile)
}

func handleRegister(c echo.Context) error {
	newUser := new(RegisterUserModel)

	if err := c.Bind(newUser); err != nil {
		return err
	}

	if err := c.Validate(newUser); err != nil {
		return err
	}

	if validPassword := utils.ValidatePassword(newUser.Password); !validPassword {
		return types.NewPasswordError()
	}

	/* hash, err := utils.HashPassword(newUser.Password)

	if err != nil {
		return err
	}

	newUser.Password = hash */

	userRecord, err := CreateUser(context.Background(), newUser.UserName, newUser.Email, newUser.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.Map{"Success": "User created", "User": userRecord})
}

func handleLogin(c echo.Context) error {

	user := new(LoginUserModel)

	if err := c.Bind(user); err != nil {
		return err
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	if validPassword := utils.ValidatePassword(user.Password); !validPassword {
		return types.NewPasswordError()
	}

	token, err := LoginWithEmailAndPassword(user.Email, user.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, types.Map{"Status": "Loged in", "Token": token})
}

func handleGetProfile(c echo.Context) error {

	authHeader := c.Request().Header["Authorization"][0]

	if authHeader == "" {
		return types.NewApiError(http.StatusUnauthorized, "no token provided")
	}

	if !strings.HasPrefix(authHeader, "Bearer") {
		return types.NewApiError(http.StatusUnauthorized, "Invalid token format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	authToken, err := VerifyToken(context.Background(), token)
	if err != nil {
		return types.NewApiError(http.StatusUnauthorized, fmt.Sprintf("error verifying token: %v", err))
	}

	userData, err := GetUSerInfo(context.Background(), authToken.UID)
	if err != nil {
		return types.NewApiError(http.StatusInternalServerError, fmt.Sprintf("user not found: %v", err))
	}

	return c.JSON(http.StatusAccepted, types.Map{"Username": userData.DisplayName, "Email": userData.Email})
}
