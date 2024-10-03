package product

import (
	"context"
	"net/http"

	"github.com/ADahjer/egocomerce/types"
	"github.com/ADahjer/egocomerce/utils"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Group) {
	router.GET("", handleGetAll)
	router.GET("/:id", handleGetOne)
	router.GET("/category/:id", handleGetByCategorie)
	router.POST("", handleCreate)
	router.DELETE("/:id", handleDelete)
	router.PUT("/:id", handleUpdate)
}

func handleGetAll(c echo.Context) error {
	products, err := GetAllProducts(context.Background())
	if err != nil {
		return err
	}

	if products == nil {
		return c.JSON(http.StatusOK, []string{})
	}

	return c.JSON(http.StatusOK, products)

}

func handleGetOne(c echo.Context) error {
	id := c.Param("id")

	prod, err := GetProductById(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, prod)
}

func handleCreate(c echo.Context) error {
	newProd := new(CreateProductModel)
	file, err := c.FormFile("image")

	msg, validFile := utils.ValidateImageType(file)
	if !validFile {
		return types.NewApiError(http.StatusBadRequest, msg)
	}

	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	newProd.Image = src

	if err := c.Bind(newProd); err != nil {
		return err
	}

	if err := c.Validate(newProd); err != nil {
		return err
	}

	ref, err := CreateProduct(context.Background(), *newProd, src)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.Map{"Message": "Product created", "id": ref})
}

func handleDelete(c echo.Context) error {
	return nil
}

func handleUpdate(c echo.Context) error {
	return nil
}

func handleGetByCategorie(c echo.Context) error {
	return nil
}
