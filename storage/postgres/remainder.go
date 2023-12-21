package postgres

import (
	"context"
	"database/sql"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type remainderRepo struct {
	db *pgxpool.Pool
}

func NewRemainderRepo(db *pgxpool.Pool) *remainderRepo {
	return &remainderRepo{
		db: db,
	}
}

func (r *remainderRepo) Create(ctx context.Context, req *models.Remainder) (*models.Remainder, error) {

	var (
		remainderId = uuid.New().String()
		query       = `
		INSERT INTO "remainder"(
		  "id",
		  "product_id",
		  "name",
		  "branch_id",
		  "quantity",
		  "sale_price",
		  "coming_price",
		  "updated_at"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		remainderId,
		req.ProductID,
		req.Name,
		req.BranchID,
		req.Quantity,
		req.SalePrice,
		req.ComingPrice,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.RemainderPrimaryKey{Id: remainderId})
}

func (r *remainderRepo) GetByID(ctx context.Context, req *models.RemainderPrimaryKey) (*models.Remainder, error) {

	var (
		query = `
			SELECT
				 "id",
				 "product_id",
				 "branch_id",
				 "name",
				 "created_at",
				 "updated_at"
			FROM "remainder"
			WHERE id = $1
		`
	)

	var (
		Id          sql.NullString
		ProductID   sql.NullString
		ProductName sql.NullString
		BranchID    sql.NullString
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&ProductID,
		&ProductName,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Remainder{
		Id:        Id.String,
		ProductID: ProductID.String,
		BranchID:  BranchID.String,
		Name:      ProductName.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (r *remainderRepo) GetList(ctx context.Context, req *models.GetListRemainderRequest) (*models.GetListRemainderResponse, error) {
	var (
		resp  models.GetListRemainderResponse
		where = " WHERE TRUE"
	)

	if len(req.Search) > 0 {
		where += " AND (name ILIKE '%" + req.Search + "%' OR branch_id ILIKE '%" + req.Search + "%')"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
	  SELECT
		COUNT(*) OVER(),
			"id",
			"name",
			"product_id",
			"quantity",
			"coming_price",
			"sale_price",
			"branch_id",
			"created_at",
			"updated_at"
	  FROM "remainder"
	`

	query += where
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			Id          sql.NullString
			Product     sql.NullString
			Name        sql.NullString
			Quantity    sql.NullInt64
			PriceIncome sql.NullFloat64
			PriceSales  sql.NullFloat64
			BranchID    sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&Product,
			&Name,
			&Quantity,
			&PriceIncome,
			&PriceSales,
			&BranchID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Remainders = append(resp.Remainders, &models.Remainder{
			Id:          Id.String,
			ProductID:   Product.String,
			BranchID:    BranchID.String,
			Name:        Name.String,
			Quantity:    int(Quantity.Int64),
			ComingPrice: PriceIncome.Float64,
			SalePrice:   PriceSales.Float64,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *remainderRepo) Update(ctx context.Context, req *models.Remainder) (int64, error) {

	query := `
		UPDATE "remainder"
			SET
				"product_id" = $2,
				"quantity" = $3,
				"coming_price" = $4,
				"sale_price" = $5,
				"branch_id" = $6,
				"name" = $7,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.ProductID,
		req.Quantity,
		req.ComingPrice,
		req.SalePrice,
		req.BranchID,
		req.Name,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *remainderRepo) Delete(ctx context.Context, req *models.RemainderPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM remainder WHERE id = $1", req.Id)
	return err
}
