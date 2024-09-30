package main

import (
	"github.com/ADahjer/egocomerce/pkg/user"
	"github.com/ADahjer/egocomerce/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.HideBanner = true

	loggerConfig := &middleware.LoggerConfig{
		Format:           "[${time_custom}]:  Method ${method}  Status ${status}  Path ${uri}\n",
		CustomTimeFormat: "1/2/2006 15:04:05",
	}

	e.Use(middleware.LoggerWithConfig(*loggerConfig))
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = utils.ApiErrorHandler

	api := e.Group("/api/v1")
	user.RegisterRoutes(api)

	e.Logger.Fatal(e.Start(":3000"))
}
