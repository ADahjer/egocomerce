package user

import (
	"context"
	"net/http"

	"github.com/ADahjer/egocomerce/types"
	"github.com/ADahjer/egocomerce/utils"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Group) {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)
	router.GET("/profile", handleGetProfile, AuthMiddleware)
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

	user := c.Get("user")

	if user == nil {
		return types.NewApiError(http.StatusUnauthorized, "error getting user data from context")
	}

	return c.JSON(200, user)

}
