package postgres

import (
	"context"
	"database/sql"
	"exam/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type clientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (r *clientRepo) Create(ctx context.Context, req *models.CreateClient) (*models.Client, error) {

	var (
		clientId = uuid.New().String()
		query    = `
			INSERT INTO "client"(
				"id",
				"first_name",
				"last_name",
				"father_name",
				"phone",
				"birthday",
				"gender",
				"branch_id",
				"active",
				"updated_at"
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, NOW())`
		birthday, err = time.Parse("2006-01-02", req.Birthday)
	)
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx,
		query,
		clientId,
		req.FirstName,
		req.LastName,
		req.FatherName,
		req.Phone,
		birthday.Format("2006-01-02"),
		req.Gender,
		req.BranchID,
		req.Active,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, &models.ClientPrimaryKey{Id: clientId})
}

func (r *clientRepo) GetByID(ctx context.Context, req *models.ClientPrimaryKey) (*models.Client, error) {

	var (
		query = `
			SELECT
				"id",
				"first_name",
				"last_name",
				"father_name",
				"phone",
				"birthday",
				"gender",
				"branch_id",
				"active",
				"created_at",
				"updated_at"
			FROM "client"
			WHERE "id" = $1
		`
	)

	var (
		Id         sql.NullString
		FirstName  sql.NullString
		LastName   sql.NullString
		FatherName sql.NullString
		Phone      sql.NullString
		Birthday   sql.NullString
		Gender     sql.NullString
		BranchID   sql.NullString
		Active     sql.NullBool
		CreatedAt  sql.NullString
		UpdatedAt  sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&FirstName,
		&LastName,
		&FatherName,
		&Phone,
		&Birthday,
		&Gender,
		&BranchID,
		&Active,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &models.Client{
		Id:         Id.String,
		FirstName:  FirstName.String,
		LastName:   LastName.String,
		FatherName: FatherName.String,
		Phone:      Phone.String,
		Birthday:   Birthday.String,
		Gender:     Gender.String,
		BranchID:   BranchID.String,
		Active:     Active.Bool,
		CreatedAt:  CreatedAt.String,
		UpdatedAt:  UpdatedAt.String,
	}, nil
}

func (r *clientRepo) GetList(ctx context.Context, req *models.GetListClientRequest) (*models.GetListClientResponse, error) {
	var (
		resp   models.GetListClientResponse
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
		where += fmt.Sprintf(" AND (first_name ILIKE %s OR last_name ILIKE %s OR father_name ILIKE %s OR phone ILIKE %s OR branch_id::text ILIKE %s)",
			searchParam, searchParam, searchParam, searchParam, searchParam)
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"first_name",
			"last_name",
			"father_name",
			"phone",
			"birthday",
			"gender",
			"branch_id",
			"active",
			"updated_at",
			"created_at"
		FROM "client"
	`

	query += where + sort + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id         sql.NullString
			FirstName  sql.NullString
			LastName   sql.NullString
			FatherName sql.NullString
			Phone      sql.NullString
			Birthday   sql.NullString
			Gender     sql.NullString
			BranchID   sql.NullString
			Active     sql.NullBool
			CreatedAt  sql.NullString
			UpdatedAt  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Id,
			&FirstName,
			&LastName,
			&FatherName,
			&Phone,
			&Birthday,
			&Gender,
			&BranchID,
			&Active,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Clients = append(resp.Clients, &models.Client{
			Id:         Id.String,
			FirstName:  FirstName.String,
			LastName:   LastName.String,
			FatherName: FatherName.String,
			Phone:      Phone.String,
			Birthday:   Birthday.String,
			Gender:     Gender.String,
			BranchID:   BranchID.String,
			Active:     Active.Bool,
			CreatedAt:  CreatedAt.String,
			UpdatedAt:  UpdatedAt.String,
		})
	}
	return &resp, nil
}

func (r *clientRepo) Update(ctx context.Context, req *models.UpdateClient) (int64, error) {

	query := `
		UPDATE client
			SET
				first_name = $2,
				last_name = $3,
				father_name = $4,
				phone = $5,
				birthday = $6,
				gender = $7,
				branch_id = $8,
				active = $9,
				updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Exec(
		ctx,
		query,
		req.Id,
		req.FirstName,
		req.LastName,
		req.FatherName,
		req.Phone,
		req.Birthday,
		req.Gender,
		req.BranchID,
		req.Active,
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

func (r *clientRepo) Delete(ctx context.Context, req *models.ClientPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM client WHERE id = $1", req.Id)
	return err
}
