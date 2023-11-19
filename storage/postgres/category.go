package postgres

import (
	"database/sql"
	"exam/models"
	"exam/pkg/helpers"

	"fmt"
	"log"

	"github.com/spf13/cast"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}
func (r *categoryRepo) Create(req *models.CreateCategory) (*models.Category, error) {
	var (
		category models.Category
		parentId sql.NullString
		query    = `
    INSERT INTO "category"(
      "title",
      "image",
      "parent_id",
      "updated_at"
    ) VALUES ($1 , $2 , $3 , NOW()) RETURNING *`
	)

	resp := r.db.QueryRow(
		query,
		req.Title,
		req.Image,
		helpers.NewNullString(req.ParentId),
	)

	err := resp.Scan(
		&category.Id,
		&category.Title,
		&category.Image,
		&parentId,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	category.ParentId = cast.ToString(parentId)

	if err != nil {
		return nil, err
	}
	fmt.Println(category)
	return r.GetById(&models.CategoryPrimaryKey{Id: category.Id})
}

func (c *categoryRepo) GetById(req *models.CategoryPrimaryKey) (*models.Category, error) {

	var (
		parentId sql.NullString
		category models.Category
		query    = `
	SELECT 
		"id",
		"title",
		"image",
		"parent_id",
		"created_at",
		"updated_at"

	FROM "category" 
		WHERE "id" = $1
		`
	)

	err := c.db.QueryRow(query, req.Id).Scan(
		&category.Id,
		&category.Title,
		&category.Image,
		&parentId,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	category.ParentId = parentId.String

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *categoryRepo) GetList(req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {
	var (
		resp   models.GetListCategoryResponse
		where  = "WHERE TRUE "
		offset = "OFFSET 0 "
		limit  = "LIMIT 10 "
		sort   = "ORDER BY created_at DESC "
	)
	if req.Offset > 0 {
		offset = fmt.Sprintf("Offset %d ", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" Limit %d", req.Limit)
	}

	if len(req.Search) > 0 {
		where += "AND title ILIKE " + " '%" + req.Search + "%' "
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
	SELECT 
		COUNT(*) OVER(), 
		"id",
		"title",
		"image",
		"parent_id",
		"created_at",
		"updated_at" 
	FROM "category"
	`

	query += where + sort + offset + limit

	fmt.Println(query)
	rows, err := c.db.Query(query)
	if err != nil {
		fmt.Println(query)
		return nil, err
	}

	for rows.Next() {

		var (
			category models.Category
			parentID sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&category.Id,
			&category.Title,
			&category.Image,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		category.ParentId = parentID.String
		resp.Categories = append(resp.Categories, &category)
		fmt.Println(category)
	}

	return &resp, nil
}

func (c *categoryRepo) Delete(req *models.CategoryPrimaryKey) (string, error) {
	_, err := c.db.Exec("DELETE FROM category WHERE id = $1 ", req.Id)
	if err != nil {
		log.Println(err.Error())
	}

	return "Deleted", nil
}

func (c *categoryRepo) Update(req *models.UpdateCategory) (int64, error) {

	query := `
	UPDATE category 
		SET 
			title = $2,
			image = $3 ,
			parent_id = $4
	WHERE id = $1 
	`

	result, err := c.db.Exec(
		query,
		req.Id,
		req.Title,
		req.Image,
		helpers.NewNullString(req.ParentId),
	)
	if err != nil {
		return 0, nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
