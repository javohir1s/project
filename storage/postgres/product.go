package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(ctx context.Context, req *models.CreateProduct) (*models.Product, error) {

	var (
		productId = uuid.New().String()
		query     = `
			INSERT INTO "product"(
				"id",
				"name",
				"price",
				"branch_id",
				"updated_at"
			) VALUES ($1, $2, $3,$4, NOW())`
	)
	_, err := r.db.Exec(ctx,
		query,
		productId,
		req.Name,
		req.Price,
		req.BranchID,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.ProductPrimaryKey{Id: productId})
}

func (r *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query = `
			SELECT
				"id",
				"name",
				"price",
				"branch_id",
				"created_at",
				"updated_at"
			FROM "product"
			WHERE id = $1
		`
	)

	var (
		Id        sql.NullString
		Name      sql.NullString
		Price     sql.NullFloat64
		BranchID  sql.NullString
		CreatedAt sql.NullString
		UpdatedAt sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&Name,
		&Price,
		&BranchID,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:        Id.String,
		Name:      Name.String,
		Price:     Price.Float64,
		BranchID:  BranchID.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (r *productRepo) GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error) {
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
		searchParam := fmt.Sprintf("'%s%%'", req.Search)
		where += fmt.Sprintf(" AND (name ILIKE %s OR  branch_id::text ILIKE %s)",
			searchParam, searchParam)
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"price",
			"branch_id",
			"created_at",
			"updated_at"
		FROM "product"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id        sql.NullString
			Name      sql.NullString
			Price     sql.NullFloat64
			BranchID  sql.NullString
			CreatedAt sql.NullString
			UpdatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&Name,
			&Price,
			&BranchID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Products = append(resp.Products, &models.Product{
			Id:        Id.String,
			Name:      Name.String,
			Price:     Price.Float64,
			BranchID:  BranchID.String,
			CreatedAt: CreatedAt.String,
			UpdatedAt: UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *productRepo) Update(ctx context.Context, req *models.UpdateProduct) (int64, error) {

	query := `
		UPDATE "product"
			SET
				"name" = $2,
				"price" = $3,
				"branch_id" = $4,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Name,
		req.Price,
		req.BranchID,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)
	return err
}
