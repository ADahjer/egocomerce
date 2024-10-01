package category

type CreateCategoryModel struct {
	Name string `json:"name" validate:"required,min=5"`
}
