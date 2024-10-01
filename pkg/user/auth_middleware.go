package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ADahjer/egocomerce/types"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

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

		c.Set("user", types.Map{"uid": userData.UID, "username": userData.DisplayName, "email": userData.Email, "claims": userData.CustomClaims})

		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(types.Map)

		claims := user["claims"].(map[string]interface{})

		isAdmin, ok := claims["admin"].(bool)
		if !ok || !isAdmin {
			return types.NewApiError(http.StatusForbidden, "not enogh privileges")
		}

		return next(c)
	}
}
