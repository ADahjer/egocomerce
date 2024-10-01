package main

import (
	"github.com/ADahjer/egocomerce/database"
	"github.com/ADahjer/egocomerce/pkg/category"
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

	e.HTTPErrorHandler = utils.ApiErrorHandler
	e.Validator = &types.CustomValidator{Validator: validator.New()}

	database.NewStore()
	user.InitRepo()
	category.InitRepo()

	api := e.Group("/api/v1")
	user.RegisterRoutes(api)

	categoryRouter := api.Group("/category")
	category.RegisterRoutes(categoryRouter)

	e.Logger.Fatal(e.Start(":3000"))
}
