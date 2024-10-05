package product

import "mime/multipart"

// What the API gets from the client
type CreateProductModel struct {
	Image      multipart.File `json:"image" form:"image"`
	Name       string         `json:"name" form:"name" validate:"required,min=4"`
	Price      float64        `json:"price" form:"price" validate:"required"`
	Categories []string       `json:"categories" form:"categories" validate:"required,min=1"`
	Discount   float64        `json:"discount" form:"discount"`
}

// what its saved to the DB
type InsertProductModel struct {
	Image      string      `json:"image" form:"image"`
	Name       string      `json:"name" form:"name" validate:"required,min=4"`
	Price      float64     `json:"price" form:"price" validate:"required"`
	Categories interface{} `json:"categories" form:"categories" validate:"required,min=1"`
	Discount   float64     `json:"discount" form:"discount"`
}

// what we get from the DB
type ProductModel struct {
	Image      string      `json:"image"`
	Id         string      `json:"id"`
	Name       string      `json:"name" validate:"required,min=4"`
	Price      float64     `json:"price" validate:"required"`
	Categories interface{} `json:"categories" validate:"required,min=1"`
	Discount   float64     `json:"discount"`
}
