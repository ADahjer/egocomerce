package category

import (
	"context"
	"net/http"

	"github.com/ADahjer/egocomerce/types"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Group) {
	router.POST("/", handleCreate)
	router.GET("/", handleGetAll)
	router.GET("/:id", handleGetById)
	router.DELETE("/:id", handleDelete)
	router.PUT("/:id", handleUpdate)
}

func handleCreate(c echo.Context) error {
	newCtg := new(CreateCategoryModel)

	if err := c.Bind(newCtg); err != nil {
		return err
	}

	if err := c.Validate(newCtg); err != nil {
		return err
	}

	docRef, err := CreateCategory(context.Background(), newCtg.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.Map{"Message": "New category created", "NewCategoryId": docRef})
}

func handleGetAll(c echo.Context) error {
	categories, err := GetAllCategories(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func handleGetById(c echo.Context) error {
	id := c.Param("id")

	ctg, err := GetCategoryById(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ctg)
}

func handleDelete(c echo.Context) error {
	id := c.Param("id")

	res, err := DeleteCategory(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.Map{"Delete": res})
}

func handleUpdate(c echo.Context) error {

	newCtg := new(CreateCategoryModel)

	if err := c.Bind(newCtg); err != nil {
		return err
	}

	if err := c.Validate(newCtg); err != nil {
		return err
	}

	id := c.Param("id")

	updatedCtg, err := UpdateCategory(context.Background(), id, newCtg.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.Map{"Message": "Category updated", "UpdatedCategory": updatedCtg})
}
