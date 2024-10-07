package cart

import (
	"context"

	"github.com/labstack/echo/v4"
)

// checks in DB if the cart exists and if it doesnt, creates a new one
func ActiveCartMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return err
		}

		cart, cartId, err := getActiveCart(context.Background(), userId)
		if err != nil {
			return err
		}

		if cart == nil && cartId == "" {
			_, _, err = CreateNewCart(context.Background(), userId)
			if err != nil {
				return err
			}
		}

		return next(c)
	}
}
