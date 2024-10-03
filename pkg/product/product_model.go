package product

import "mime/multipart"

type CreateProductModel struct {
	Image      multipart.File `json:"image" form:"image" validate:"required"`
	Name       string         `json:"name" form:"name" validate:"required,min=4"`
	Price      float64        `json:"price" form:"price" validate:"required"`
	Categories []string       `json:"categories" form:"categories" validate:"required,min=1"`
}

type ProductModel struct {
	Image      string      `json:"image"`
	Id         string      `json:"id"`
	Name       string      `json:"name" validate:"required,min=4"`
	Price      float64     `json:"price" validate:"required"`
	Categories interface{} `json:"categories" validate:"required,min=1"`
}
