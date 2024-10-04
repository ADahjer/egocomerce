package cart

import (
	"context"
	"net/http"

	"github.com/ADahjer/egocomerce/pkg/user"
	"github.com/ADahjer/egocomerce/types"
	"github.com/labstack/echo/v4"
)

// Aux func
func getUserId(c echo.Context) (string, error) {
	userI := c.Get("user")

	if userI == nil {
		return "", types.NewApiError(http.StatusUnauthorized, "error getting user data")
	}

	user := userI.(types.Map)
	id := user["uid"].(string)

	return id, nil
}

func RegisterRoutes(router *echo.Group) {
	router.GET("", handleGetActiveCart, user.AuthMiddleware)
	router.GET("/:id", handleGetOne)
	router.POST("", handleCreate, user.AuthMiddleware)
}

func handleCreate(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	newItem := new(CartItemModel)
	if err := c.Bind(newItem); err != nil {
		return err
	}

	if err := c.Validate(newItem); err != nil {
		return err
	}

	err = AddItemToCart(context.Background(), userId, *newItem)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.Map{"created": true})
}

// usefull to get previous carts (or completed)
func handleGetOne(c echo.Context) error {
	return nil
}

func handleGetActiveCart(c echo.Context) error {
	userId, err := getUserId(c)

	if err != nil {
		return err
	}

	cart, cartId, err := getActiveCart(context.Background(), userId)
	if err != nil {
		return err
	}

	if cart == nil && cartId == "" {
		cart, cartId, err = CreateNewCart(context.Background(), userId)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, types.Map{"active_cart_id": cartId, "active_cart_data": cart})

}
