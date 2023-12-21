package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ComingRepo struct {
	db *pgxpool.Pool
}

func NewComingRepo(db *pgxpool.Pool) *ComingRepo {
	return &ComingRepo{
		db: db,
	}
}

func (r *ComingRepo) Create(ctx context.Context, req *models.CreateComing) (*models.Coming, error) {

	var (
		comingId = uuid.New().String()
		query    = `
			INSERT INTO "coming"(
				"id",
				"increment_id",
				"branch_id",
				"updated_at"
			) VALUES ($1, $2, $3, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		comingId,
		req.IncrementID,
		req.BranchID,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.ComingPrimaryKey{Id: comingId})
}

func (r *ComingRepo) GetByID(ctx context.Context, req *models.ComingPrimaryKey) (*models.Coming, error) {

	var (
		query = `
			SELECT
				 "id",
				 "increment_id",
				 "branch_id",
				 "created_at",
				 "updated_at"
			FROM "coming"
			WHERE id = $1
		`
	)

	var (
		Id          sql.NullString
		IncrementID sql.NullString
		BranchID    sql.NullString
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&BranchID,
		&IncrementID,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Coming{
		Id:          Id.String,
		IncrementID: IncrementID.String,
		BranchID:    BranchID.String,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (r *ComingRepo) GetList(ctx context.Context, req *models.GetListComingRequest) (*models.GetListComingResponse, error) {
	var (
		resp   models.GetListComingResponse
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
			"increment_id",
			"branch_id",
			"created_at",
			"updated_at"
		FROM "coming"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id          sql.NullString
			IncrementID sql.NullString
			BranchID    sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&IncrementID,
			&BranchID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Cominges = append(resp.Cominges, &models.Coming{
			Id:          Id.String,
			IncrementID: IncrementID.String,
			BranchID:    BranchID.String,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *ComingRepo) Update(ctx context.Context, req *models.UpdateComing) (int64, error) {

	query := `
		UPDATE "coming"
			SET
				"branch_id" = $2,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.BranchID,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *ComingRepo) Delete(ctx context.Context, req *models.ComingPrimaryKey) error {
	query := `DELETE FROM coming WHERE id = $1`
	_, err := r.db.Exec(ctx, query, req.Id)
	return err
}
