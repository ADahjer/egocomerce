package product

type CreateProductModel struct {
	Name       string   `json:"name" validate:"required,min=4"`
	Price      float32  `json:"price" validate:"required"`
	Categories []string `json:"categories" validate:"required,min=1"`
}

type ProductModel struct {
	Id         string      `json:"id"`
	Name       string      `json:"name" validate:"required,min=4"`
	Price      float64     `json:"price" validate:"required"`
	Categories interface{} `json:"categories" validate:"required,min=1"`
}
