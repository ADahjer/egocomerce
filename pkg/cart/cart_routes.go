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
		return "", types.NewApiError(http.StatusUnauthorized, "error getting user data, token may have expired")
	}

	user := userI.(types.Map)
	id := user["uid"].(string)

	return id, nil
}

func RegisterRoutes(router *echo.Group) {
	router.GET("", handleGetActiveCart, user.AuthMiddleware)
	router.GET("/completed", handleGetCompletedCarts, user.AuthMiddleware)
	router.POST("", handleAddItem, user.AuthMiddleware)
	router.PUT("/complete", handleComplete, user.AuthMiddleware)
	router.DELETE("", hanldeDeleteCart, user.AuthMiddleware)
}

func handleAddItem(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	newItem := new(NewCartItemModel)
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

func hanldeDeleteCart(c echo.Context) error {
	// check if the user's active cart its empty or not
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	err = VoidCart(context.Background(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.Map{"cart_void": true})
}

func handleComplete(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	err = CompleteCart(context.Background(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.Map{"cart_completed": true})
}

func handleGetCompletedCarts(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	carts, err := GetCompletedCarts(context.Background(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, carts)
}
