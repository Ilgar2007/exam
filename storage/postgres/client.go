package postgres

import (
	"database/sql"
	"exam/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type clientRepo struct {
	db *sql.DB
}

func NewClientRepo(db *sql.DB) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (r *clientRepo) Create(req *models.CreateClient) (*models.Client, error) {
	clientID := uuid.New().String()
	query := `
		INSERT INTO "clients" (
			"first_name", 
			"last_name", 
			"phone",
			"photo",
			"date_of_birth",
			"created_at"
		) VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err := r.db.Exec(
		query,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Photo,
		req.DateOfBirth,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.ClientPrimaryKey{ID: clientID})
}

func (c *clientRepo) GetByID(req *models.ClientPrimaryKey) (*models.Client, error) {

	var (
		client models.Client
		query  = `
	SELECT 
		"id",
		"first_name", 
		"last_name",
    "phone",
    "photo",
		"date_of_birth",
		"created_at"
	FROM "clients" 
		WHERE "id" = $1
		`
	)

	err := c.db.QueryRow(query, req.ID).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.Phone,
		&client.Photo,
		&client.DateOfBirth,
		&client.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *clientRepo) GetList(req *models.GetListClientRequest) (*models.GetListClientResponse, error) {
	var (
		resp   models.GetListClientResponse
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
		where += fmt.Sprintf(" AND (first_name ILIKE '%%%s%%' OR last_name ILIKE '%%%s%%' OR phone ILIKE '%%%s%%')",
			req.Search, req.Search, req.Search)
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
	SELECT 
		COUNT(*) OVER(), 
		"id",
		"first_name", 
		"last_name",
    "phone",
    "photo",
		"date_of_birth",
		"created_at"
	FROM "clients"
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
			client models.Client
		)

		err = rows.Scan(
			&resp.Count,
			&client.ID,
			&client.FirstName,
			&client.LastName,
			&client.Phone,
			&client.Photo,
			&client.DateOfBirth,
			&client.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Clients = append(resp.Clients, &client)
		fmt.Println(client)
	}

	return &resp, nil
}

func (c *clientRepo) Delete(req *models.ClientPrimaryKey) (string, error) {
	_, err := c.db.Exec("DELETE FROM clients WHERE id = $1", req.ID)
	if err != nil {
		log.Println(err.Error())
	}
	return "Deleted", err
}

func (r *clientRepo) Update(req *models.UpdateClient) (int64, error) {
	query := `
		UPDATE clients
			SET
				first_name = $2,   
				last_name = $3,
				phone = $4,   
				photo =$5 ,
				date_of_birth = $6,
				updated_at = NOW()  
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.ID,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Photo,
		req.DateOfBirth,
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
