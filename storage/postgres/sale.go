package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SaleRepo struct {
	db *pgxpool.Pool
}

func NewSaleRepo(db *pgxpool.Pool) *SaleRepo {
	return &SaleRepo{
		db: db,
	}
}

func (r *SaleRepo) Create(ctx context.Context, req *models.CreateSale) (*models.Sale, error) {

	var (
		SaleId = uuid.New().String()
		query       = `
			INSERT INTO "sale"(
				"id",
				"branch_id",
				"client_id",
				"increment_id",
				"total_price",
				"updated_at"
			) VALUES ($1, $2, $3,$4,$5, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		SaleId,
		req.BranchID,
		req.ClientID,
		req.IncrementID,
		req.TotalPrice,

	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.SalePrimaryKey{Id: SaleId})
}

func (r *SaleRepo) GetByID(ctx context.Context, req *models.SalePrimaryKey) (*models.Sale, error) {

	var (
		query = `
			SELECT
				 "id",
				 "branch_id",
				 "client_id",
				 "increment_id",
				 "total_price",
				 "paid",
				 "debt",
				 "created_at",
				 "updated_at"
			FROM "sale"
			WHERE id = $1
		`
	)

	var (
		Id          sql.NullString
		BranchID    sql.NullString
		ClientID    sql.NullString
		IncrementID sql.NullString
		TotalPrice  sql.NullFloat64
		Paid        sql.NullFloat64
		Debd        sql.NullFloat64
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&ClientID,
		&IncrementID,
		&TotalPrice,
		&Paid,
		&Debd,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Sale{
		Id:          Id.String,
		BranchID:    BranchID.String,
		ClientID:    ClientID.String,
		IncrementID: IncrementID.String,
		TotalPrice:  TotalPrice.Float64,
		Paid:        Paid.Float64,
		Debd:        Paid.Float64,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *SaleRepo) GetList(ctx context.Context, req *models.GetListSaleRequest) (*models.GetListSaleResponse, error) {
	var (
		resp   models.GetListSaleResponse
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

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"branch_id",
			"client_id",
			"increment_id",
			"total_price",
			"paid",
			"debt",
			"created_at",
			"updated_at"
		FROM "sale"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id          sql.NullString
			BranchID    sql.NullString
			ClientID    sql.NullString
			IncrementID sql.NullString
			TotalPrice  sql.NullFloat64
			Paid        sql.NullFloat64
			Debd        sql.NullFloat64
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&BranchID,
			&ClientID,
			&IncrementID,
			&TotalPrice,
			&Paid,
			&Debd,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Sales = append(resp.Sales, &models.Sale{
			Id:          Id.String,
			BranchID:    BranchID.String,
			ClientID:    ClientID.String,
			IncrementID: IncrementID.String,
			TotalPrice:  TotalPrice.Float64,
			Paid:        Paid.Float64,
			Debd:        Paid.Float64,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *SaleRepo) Update(ctx context.Context, req *models.UpdateSale) (int64, error) {

	query := `
		UPDATE "sale"
			SET
				"branch_id" = $2,
				"client_id" = $3,
				"increment_id" = $4,
				"total_price" = $5,
				"paid" = $6,
				"debt" = $7,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	// fmt.Println(req.Id,
	// 	req.BranchID,
	// 	req.ClientID,
	// 	req.IncrementID,
	// 	req.TotalPrice,
	// 	req.Paid,
	// 	req.Debd,
	// 	req.TotalPrice,)
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.BranchID,
		req.ClientID,
		req.IncrementID,
		req.TotalPrice,
		req.Paid,
		req.Debd,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *SaleRepo) Delete(ctx context.Context, req *models.SalePrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sale WHERE id = $1", req.Id)
	return err
}
