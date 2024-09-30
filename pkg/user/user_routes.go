package user

import (
	"net/http"

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

	hash, err := utils.HashPassword(newUser.Password)

	if err != nil {
		return err
	}

	newUser.Password = hash

	return c.JSON(http.StatusCreated, types.Map{"Success": "User created", "User": newUser.UserName, "hash": newUser.Password})
}

func handleLogin(c echo.Context) error {
	return c.String(200, "Login")
}

func handleGetProfile(c echo.Context) error {
	return c.String(200, "Profile")
}
