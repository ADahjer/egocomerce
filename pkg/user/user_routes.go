package user

import "github.com/labstack/echo/v4"

func RegisterRoutes(router *echo.Group) {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)
	router.GET("/profile", handleGetProfile)
}

func handleRegister(c echo.Context) error {
	return c.String(200, "Register")
}

func handleLogin(c echo.Context) error {
	return c.String(200, "Login")
}

func handleGetProfile(c echo.Context) error {
	return c.String(200, "Profile")
}
