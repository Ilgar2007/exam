package models

type Category struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	ParentId  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateCategory struct {
	Title    string `json:"title"`
	Image    string `json:"image"`
	ParentId string `json:"parent_id"`
}

type UpdateCategory struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	ParentId string `json:"parent_id"`
}

type CategoryPrimaryKey struct {
	Id string `json:"id"`
}

type GetListCategoryRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListCategoryResponse struct {
	Count      int         `json:"count"`
	Categories []*Category `json:"categories"`
}
