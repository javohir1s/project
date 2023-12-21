package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pickingListRepo struct {
	db *pgxpool.Pool
}

func NewPickingListRepo(db *pgxpool.Pool) *pickingListRepo {
	return &pickingListRepo{
		db: db,
	}
}

func (r *pickingListRepo) Create(ctx context.Context, req *models.PickingList) (*models.PickingList, error) {

	var (
		pickingListId = uuid.New().String()
		query         = `
			INSERT INTO "picking_list"(
				"id",
				"product_id",
				"quantity",
				"coming_id",
				"coming_increment_id",
				"updated_at"
			) VALUES ($1, $2, $3, $4, $5, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		pickingListId,
		req.Product_ID,
		req.Quantity,
		req.ComingID,
		req.ComingIncrementID,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.PickingListPrimaryKey{Id: pickingListId})
}

func (r *pickingListRepo) GetByID(ctx context.Context, req *models.PickingListPrimaryKey) (*models.PickingList, error) {

	var (
		query = `
			SELECT
				"id",
				"product_id",
				"price",
				"quantity",
				"total_price",
				"coming_id",
				"coming_increment_id",
				"created_at",
				"updated_at"
			FROM "picking_list"
			WHERE id = $1
		`
	)

	var (
		ID                sql.NullString
		Product_ID        sql.NullString
		Price             sql.NullFloat64
		Quantity          sql.NullInt64
		Total_price       sql.NullFloat64
		ComingID          sql.NullString
		ComingIncrementID sql.NullString
		CreatedAt         sql.NullString
		UpdatedAt         sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&ID,
		&Product_ID,
		&Price,
		&Quantity,
		&Total_price,
		&ComingID,
		&ComingIncrementID,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.PickingList{
		ID:                ID.String,
		Product_ID:        Product_ID.String,
		Price:             Price.Float64,
		Quantity:          int(Quantity.Int64),
		Total_price:       Total_price.Float64,
		ComingID:          ComingID.String,
		ComingIncrementID: ComingIncrementID.String,
		CreatedAt:         CreatedAt.String,
		UpdatedAt:         UpdatedAt.String,
	}, nil
}

func (r *pickingListRepo) GetList(ctx context.Context, req *models.GetListPickingListRequest) (*models.GetListPickingListResponse, error) {
	var (
		resp   models.GetListPickingListResponse
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
				"product_id",
				"price",
				"quantity",
				"total_price",
				"coming_id",
				"coming_increment_id",
				"created_at",
				"updated_at"
		FROM "picking_list"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID                sql.NullString
			Product_ID        sql.NullString
			Price             sql.NullFloat64
			Quantity          sql.NullInt64
			Total_price       sql.NullFloat64
			ComingID          sql.NullString
			ComingIncrementID sql.NullString
			CreatedAt         sql.NullString
			UpdatedAt         sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&ID,
			&Product_ID,
			&Price,
			&Quantity,
			&Total_price,
			&ComingID,
			&ComingIncrementID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Pickinges = append(resp.Pickinges, &models.PickingList{
			ID:                ID.String,
			Product_ID:        Product_ID.String,
			Price:             Price.Float64,
			Quantity:          int(Quantity.Int64),
			Total_price:       Total_price.Float64,
			ComingID:          ComingID.String,
			ComingIncrementID: ComingIncrementID.String,
			CreatedAt:         CreatedAt.String,
			UpdatedAt:         UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *pickingListRepo) Update(ctx context.Context, req *models.PickingList) (int64, error) {

	query := `
		UPDATE "picking_list"
			SET
				"product_id" = $2,
				"quantity" = $3,
				"coming_id" = $4,
				"coming_increment_id" = $5,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.ID,
		req.Product_ID,
		req.Quantity,
		req.ComingID,
		req.ComingIncrementID,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *pickingListRepo) Delete(ctx context.Context, req *models.PickingListPrimaryKey) error {
	query := `
		DELETE FROM 
			"picking_list" 
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, req.Id)
	return err
}