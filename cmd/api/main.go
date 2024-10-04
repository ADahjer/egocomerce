package main

import (
	"os"

	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/pkg/cart"
	"github.com/ADahjer/egocomerce/pkg/product"
	"github.com/ADahjer/egocomerce/pkg/user"
	"github.com/ADahjer/egocomerce/types"
	"github.com/ADahjer/egocomerce/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()

	e := echo.New()
	e.HideBanner = true

	loggerConfig := &middleware.LoggerConfig{
		Format:           "[${time_custom}]:  Method ${method}  Status ${status}  Path ${uri}\n",
		CustomTimeFormat: "1/2/2006 15:04:05",
	}

	e.Use(middleware.LoggerWithConfig(*loggerConfig))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.HTTPErrorHandler = utils.ApiErrorHandler
	e.Validator = &types.CustomValidator{Validator: validator.New()}

	database.NewStore()
	user.InitRepo()
	product.InitRepo()
	cart.InitRepo()

	api := e.Group("/v1")
	user.RegisterRoutes(api)

	productRouter := api.Group("/product")
	product.RegisterRoutes(productRouter)

	cartRouter := api.Group("/cart")
	cart.RegisterRoutes(cartRouter)

	port := os.Getenv("API_PORT")

	if port == "" {
		port = "5000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
