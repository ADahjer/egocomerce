package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.HideBanner = true

	loggerConfig := &middleware.LoggerConfig{
		Format:           "[${time_custom}]:  Path ${uri}  Method ${method}  Status ${status}",
		CustomTimeFormat: "1/2/2006 15:04:05",
	}

	e.Use(middleware.LoggerWithConfig(*loggerConfig))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":3000"))

}
