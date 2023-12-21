package postgres

import (
	"context"
	"database/sql"
	"exam/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ReportRepo struct {
	db *pgxpool.Pool
}

func NewReportRepo(db *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{
		db: db,
	}
}

func (r *ReportRepo) GetListReport(ctx context.Context, fromDate, toDate time.Time) (*models.GetListClientResponse, error) {
	var (
		respons models.GetListClientResponse
		query   = `
			SELECT 
				COUNT(*) OVER(),
				"id",
				"first_name",
				"last_name",
				"father_name",
				"phone",
				"birthday",
				"active",
				"gender",
				"branch_id",
				"created_at",
				"updated_at"
			FROM client
			WHERE "created_at" >= $1 AND "created_at" <= $2
			ORDER BY created_at DESC
		`
	)

	rows, err := r.db.Query(ctx, query, fromDate, toDate)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			Id          sql.NullString
			FirstName   sql.NullString
			LastName    sql.NullString
			FatherName  sql.NullString
			PhoneNumber sql.NullString
			BirthDay    sql.NullString
			IsActive    sql.NullBool
			Gender      sql.NullString
			BranchID    sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err := rows.Scan(
			&respons.Count,
			&Id,
			&FirstName,
			&LastName,
			&FatherName,
			&PhoneNumber,
			&BirthDay,
			&IsActive,
			&Gender,
			&BranchID,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		respons.Clients = append(respons.Clients, &models.Client{
			Id:         Id.String,
			FirstName:  FirstName.String,
			LastName:   LastName.String,
			FatherName: FatherName.String,
			Phone:      PhoneNumber.String,
			Birthday:   BirthDay.String,
			Active:     IsActive.Bool,
			Gender:     Gender.String,
			BranchID:   BranchID.String,
			CreatedAt:  CreatedAt.String,
			UpdatedAt:  UpdatedAt.String,
		})
	}

	return &respons, nil
}

func (r *ReportRepo) GetListSaleBranch(ctx context.Context, fromDate, toDate time.Time) (*models.AllSaleReport, error) {
	var (
		respons models.AllSaleReport
		query   = `
			SELECT 
				B.name,
				SUM(SP.quantity) AS quantity,
				SUM("price") AS price
			FROM "sale_product" AS SP
			JOIN "sale" AS S ON S.id = SP.sale_id
			JOIN "branch" AS B ON B.id = S.branch_id
			WHERE S.created_at >= $1 AND S.created_at <= $2
			GROUP BY B.name;
		`
	)

	rows, err := r.db.Query(ctx, query, fromDate, toDate)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			Name     sql.NullString
			Quantity sql.NullInt64
			Price    sql.NullFloat64
		)

		err := rows.Scan(
			&Name,
			&Quantity,
			&Price,
		)
		if err != nil {
			return nil, err
		}

		respons.SaleReports = append(respons.SaleReports, &models.SaleReport{
			Name:     Name.String,
			Quantity: Quantity.Int64,
			Price:    Price.Float64,
		})
	}

	return &respons, nil
}
