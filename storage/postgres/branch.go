package postgres

import (
	"database/sql"
	"exam/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type branchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(req *models.CreateBranch) (*models.Branch, error) {
	branchID := uuid.New().String()
	query := `
		INSERT INTO "branches" (
			"name", 
			"phone",
			"photo",
			"work_start_hour", 
			"work_end_hour",
			"address",
			"delivery_price",
			"active",
			"created_at"
		) VALUES ($1, $2, $3, $4, $5, $6, $7 , $8 , NOW())`

	_, err := r.db.Exec(
		query,
		req.Name,
		req.Phone,
		req.Photo,
		req.WorkStartHour,
		req.WorkEndHour,
		req.Address,
		req.DeliveryPrice,
		req.Active,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.BranchPrimaryKey{ID: branchID})
}

func (r *branchRepo) GetByID(req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		query = `
    SELECT
			"id",
			"name", 
			"phone",
			"photo",
			"work_start_hour", 
			"work_end_hour",
			"address",
			"delivery_price",
			"active",
			"created_at",
      "updated_at"
    FROM "branches"
    WHERE "id" = $1
    `
	)

	var (
		id            sql.NullString
		name          sql.NullString
		phone         sql.NullString
		photo         sql.NullString
		workStartHour sql.NullString
		workEndHour   sql.NullString
		address       sql.NullString
		deliveryPrice sql.NullFloat64
		active        sql.NullBool
		created_at    sql.NullString
		updated_at    sql.NullString
	)

	err := r.db.QueryRow(query, req.ID).Scan(
		&id,
		&name,
		&phone,
		&photo,
		&workStartHour,
		&workEndHour,
		&address,
		&deliveryPrice,
		&active,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.Branch{
		ID:            id.String,
		Name:          name.String,
		Phone:         phone.String,
		Photo:         photo.String,
		WorkStartHour: workStartHour.String,
		WorkEndHour:   workEndHour.String,
		Address:       address.String,
		DeliveryPrice: deliveryPrice.Float64,
		Active:        active.Bool,
		CreatedAt:     created_at.String,
		UpdatedAt:     updated_at.String,
	}, nil
}

func (r *branchRepo) Update(req *models.UpdateBranch) (int64, error) {
	query := `
		UPDATE branches
			SET
				name = $2,   
				phone = $3,   
				photo =$4 ,
				work_start_hour = $5,
				work_end_hour = $6,
				address= $7 ,
				delivery_price= $8,
				active= $9,
				updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.ID,
		req.Name,
		req.Phone,
		req.Photo,
		req.WorkStartHour,
		req.WorkEndHour,
		req.Address,
		req.DeliveryPrice,
		req.Active,
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

func (b *branchRepo) Delete(req *models.BranchPrimaryKey) (string, error) {
	_, err := b.db.Exec(`DELETE FROM branches WHERE id = $1`, req.ID)
	if err != nil {
		log.Println(err.Error())
	}
	return "Deleted", nil
}
func (r *branchRepo) GetList(req *models.GetListBranchRequest) (*models.GetListBranchResponse, error) {
	var (
		resp   models.GetListBranchResponse
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
		where += " AND name ILIKE" + " '%" + req.Search + "%'"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	if len(req.Name) > 0 {
		where += " AND name ILIKE '%" + req.Name + "%'"
	}
	if len(req.FromDate) > 0 || len(req.ToDate) > 0 {
		where += " AND created_at BETWEEN '" + req.FromDate + "' AND '" + req.ToDate + "'"
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"phone",
			"photo",
			"work_start_hour",
			"work_end_hour",
			"address",
			"delivery_price",
			"active",
			"created_at"
		FROM "branches"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branch models.Branch

		err = rows.Scan(
			&resp.Count,
			&branch.ID,
			&branch.Name,
			&branch.Phone,
			&branch.Photo,
			&branch.WorkStartHour,
			&branch.WorkEndHour,
			&branch.Address,
			&branch.DeliveryPrice,
			&branch.Active,
			&branch.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Branches = append(resp.Branches, &branch)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}
