package postgres

import (
	"database/sql"
	"exam/models"
	"exam/pkg/helpers"
	"fmt"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}
func (r productRepo) Create(req *models.CreateProduct) (*models.Product, error) {

	var (
		product      models.Product
		productId, _ = helpers.NewIncrementId(r.db, "products", 8)
		query        = `
    INSERT INTO "products"(
      "id",
      "title",
      "photo",
      "price",
      "description",
      "category_id",
      "updated_at"
    ) VALUES ($1 , $2 , $3 ,$4, $5, $6, NOW()) RETURNING *`
	)

	resp := r.db.QueryRow(
		query,
		productId(),
		req.Title,
		req.Photo,
		req.Price,
		req.Description,
		req.CategoryId,
	)

	err := resp.Scan(
		&product.Id,
		&product.CategoryId,
		&product.Title,
		&product.Description,
		&product.Photo,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.ProductPrimaryKey{Id: product.Id})
}

func (r *productRepo) GetByID(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		product models.Product
		query   = `
      SELECT
        "id",
        "title",
        "photo",
        "price",
        "description",
        "category_id",
        "created_at",
        "updated_at"  
      FROM "products"
      WHERE "id" = $1
    `
	)

	fmt.Println("**************************************", req.Id)
	err := r.db.QueryRow(query, req.Id).Scan(
		&product.Id,
		&product.Title,
		&product.Photo,
		&product.Price,
		&product.Description,
		&product.CategoryId,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) GetList(req *models.GetListProductRequest) (*models.GetListProductResponse, error) {
	var (
		resp   models.GetListProductResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if len(req.Search) > 0 {
		where += fmt.Sprintf(" AND (title ILIKE '%%%s%%' OR category_id ILIKE '%%%s%%' )",
			req.Search, req.Search)
	}

	var countQuery = `
		SELECT COUNT(*) FROM "products"
	`
	countQuery += where

	err := r.db.QueryRow(countQuery).Scan(&resp.Count)
	if err != nil {
		return nil, err
	}

	var selectQuery = `
		SELECT
			"id",
			"title",   
			"photo",   
			"price",
			"description",
			"category_id",
			"created_at",
			"updated_at"
		FROM "products"
	`

	selectQuery += where + sort + offset + limit
	rows, err := r.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			title       sql.NullString
			photo       sql.NullString
			price       sql.NullFloat64
			description sql.NullString
			category_id sql.NullString
			created_at  sql.NullString
			updated_at  sql.NullString
		)

		err := rows.Scan(
			&id,
			&title,
			&photo,
			&price,
			&description,
			&category_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, &models.Product{
			Id:          id.String,
			Title:       title.String,
			Photo:       photo.String,
			Price:       price.Float64,
			Description: description.String,
			CategoryId:  category_id.String,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		})
	}

	return &resp, nil
}

func (r *productRepo) Update(req *models.UpdateProduct) (int64, error) {
	query := `
		UPDATE products
			SET
				title = $2,   
				description = $3,
				photo = $4,   
				price = $5,
				category_id = $6,
				updated_at = NOW()  
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.Id,
		req.Title,
		req.Description,
		req.Photo,
		req.Price,
		helpers.NewNullString(req.CategoryId),
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *productRepo) Delete(req *models.ProductPrimaryKey) (string, error) {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", req.Id)
	return "Deleted", err
}
