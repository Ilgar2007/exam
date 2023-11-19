package models

type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateProduct struct {
	Id          string  `json:"product_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
	Photo       string  `json:"photo"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
}

type Product struct {
	Id          string  `json:"product_id"`
	Title       string  `json:"title"`
	Photo       string  `json:"photo"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateProduct struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Photo       string  `json:"photo"`
	Price       float64 `json:"price"`
	CategoryId  string  `json:"category_id"`
}

type GetListProductRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListProductResponse struct {
	Count    int        `json:"count"`
	Products []*Product `json:"products"`
}
