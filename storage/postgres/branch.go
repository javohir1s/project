package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(ctx context.Context, req *models.CreateBranch) (*models.Branch, error) {

	var (
		branchId = uuid.New().String()
		query    = `
			INSERT INTO "branch"(
				"id",
				"name", 
				"address",
				"phone",
				"updated_at"
			) VALUES ($1, $2, $3,$4, NOW())`
	)

	_, err := r.db.Exec(ctx,
		query,
		branchId,
		req.Name,
		req.Address,
		req.Phone,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.BranchPrimaryKey{Id: branchId})
}

func (r *branchRepo) GetByID(ctx context.Context, req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		query = `
			SELECT
					"id",
					"name",
					"address",
					"phone",
					"updated_at",
					"created_at"	
			FROM "branch"
			WHERE "id" = $1
		`
	)

	var (
		Id        sql.NullString
		Title     sql.NullString
		Address   sql.NullString
		Phone     sql.NullString
		CreatedAt sql.NullString
		UpdatedAt sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&Title,
		&Address,
		&Phone,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &models.Branch{
		Id:        Id.String,
		Name:      Title.String,
		Address:   Address.String,
		Phone:     Phone.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (r *branchRepo) GetList(ctx context.Context, req *models.GetListBranchRequest) (*models.GetListBranchResponse, error) {
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
		where += " AND name ILIKE" + " '%" + req.Search + "%'" + "OR phone ILIKE " + "'%" + req.Search + "%'"
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"address",
			"phone",
			"updated_at",
			"created_at"
		FROM "branch"
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
			Address   sql.NullString
			Phone     sql.NullString
			CreatedAt sql.NullString
			UpdatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&Name,
			&Address,
			&Phone,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Branches = append(resp.Branches, &models.Branch{
			Id:        Id.String,
			Name:      Name.String,
			Address:   Address.String,
			Phone:     Phone.String,
			CreatedAt: CreatedAt.String,
			UpdatedAt: UpdatedAt.String,
		})
	}
	return &resp, nil
}

func (r *branchRepo) Update(ctx context.Context, req *models.UpdateBranch) (int64, error) {

	query := `
		UPDATE branch
			SET
				name = $2,
				address = $3,
				phone = $4,
				updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Exec(
		ctx,
		query,
		req.Id,
		req.Name,
		req.Address,
		req.Phone,
	)

	if err != nil {
		return 0, err
	}
	rowsAffected := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *branchRepo) Delete(ctx context.Context, req *models.BranchPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM branch WHERE id = $1", req.Id)
	return err
}
