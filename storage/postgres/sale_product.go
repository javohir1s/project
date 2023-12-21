package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"exam/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type saleProductRepo struct {
	db *pgxpool.Pool
}

func NewSaleProductRepo(db *pgxpool.Pool) *saleProductRepo {
	return &saleProductRepo{
		db: db,
	}
}

func (r *saleProductRepo) Create(ctx context.Context, req *models.CreateSaleProduct) (*models.SaleProduct, error) {

	var (
		saleProductId = uuid.New().String()
		query         = `
			INSERT INTO "sale_product"(
				"id",
				"product_id",
				"sale_id",
				"sale_increment_id",
				"quantity",
				"price",
				"total_price"
			) VALUES ($1,$2,$3,$4,$5,$6,$7)`

		query1 = `SELECT quantity FROM remainder WHERE product_id = $1 AND branch_id = $2`
		query2 = `SELECT price FROM product WHERE id = $1`
		query3 = `UPDATE sale 
				  SET total_price = total_price + $1
				  WHERE id = $2
		`
		query4 =`SELECT branch_id from sale where id = $1`
		branchId sql.NullString
		remaining sql.NullInt64
		price sql.NullFloat64
	)
	err := r.db.QueryRow(ctx,query4,req.SaleID).Scan(&branchId,)
	if err == sql.ErrNoRows {
        return nil, errors.New("no such product")
    }
	

	err = r.db.QueryRow(ctx,query1,req.ProcutID,branchId.String).Scan(&remaining,)
	if err!=nil{
		return nil,err
	}
	fmt.Println(remaining.Int64)
	
	if remaining.Int64 < int64(req.Quantity){
		return nil, errors.New("not enough quantity")
	}

	err = r.db.QueryRow(ctx,query2,req.ProcutID).Scan(&price,)
	if err!=nil{
		return nil,err
	}


	_, err = r.db.Exec(ctx,
		query,
		saleProductId,
		req.ProcutID,
		req.SaleID,
		req.SaleIncrementID,
		req.Quantity,
		price.Float64,
		float64(req.Quantity) * price.Float64,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(float64(req.Quantity)*price.Float64)
	fmt.Println(query3)
	_,err = r.db.Exec(ctx,query3,float64(req.Quantity)*price.Float64,req.SaleID)
	if err!=nil{
		return nil,err
	}
	if remaining.Int64 > 0 || remaining.Int64 > int64(req.Quantity){
	_,err = r.db.Exec(ctx,`UPDATE remainder SET quantity = quantity - $1 where product_id = $2 AND branch_id = $3`,req.Quantity,req.ProcutID,branchId)
	if err!=nil{
		return nil,err
	}
}
	return r.GetByID(ctx, &models.SaleProductPrimaryKey{Id: saleProductId})
}

func (r *saleProductRepo) GetByID(ctx context.Context, req *models.SaleProductPrimaryKey) (*models.SaleProduct, error) {

	var (
		query = `
			SELECT
				"id",
				"product_id",
				"sale_id",
				"sale_increment_id",
				"quantity",
				"price",
				"total_price",
				"created_at",
				"updated_at"
			FROM "sale_product"
			WHERE id = $1
		`
	)

	var (
		Id              sql.NullString
		ProcutID        sql.NullString
		SaleID          sql.NullString
		SaleIncrementID sql.NullString
		Quantity        sql.NullInt64
		Price           sql.NullFloat64
		TotalPrice      sql.NullFloat64
		CreatedAt       sql.NullString
		UpdatedAt       sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&Id,
		&ProcutID,
		&SaleID,
		&SaleIncrementID,
		&Quantity,
		&Price,
		&TotalPrice,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.SaleProduct{
		Id:              Id.String,
		ProcutID:        ProcutID.String,
		SaleID:          SaleID.String,
		SaleIncrementID: SaleIncrementID.String,
		Quantity:        int(Quantity.Int64),
		Price:           Price.Float64,
		TotalPrice:      TotalPrice.Float64,
		CreatedAt:       CreatedAt.String,
		UpdatedAt:       UpdatedAt.String,
	}, nil
}

func (r *saleProductRepo) GetList(ctx context.Context, req *models.GetListSaleProductRequest) (*models.GetListSaleProductResponse, error) {
	var (
		resp   models.GetListSaleProductResponse
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
			"sale_id",
			"sale_increment_id",
			"quantity",
			"price",
			"total_price",
			"created_at",
			"updated_at"
		FROM "sale_product"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			Id              sql.NullString
			ProcutID        sql.NullString
			SaleID          sql.NullString
			SaleIncrementID sql.NullString
			Quantity        sql.NullInt64
			Price           sql.NullFloat64
			TotalPrice      sql.NullFloat64
			CreatedAt       sql.NullString
			UpdatedAt       sql.NullString
		)

		err = rows.Scan(
			&Id,
			&ProcutID,
			&SaleID,
			&SaleIncrementID,
			&Quantity,
			&Price,
			&TotalPrice,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.SaleProducts = append(resp.SaleProducts, &models.SaleProduct{
			Id:              Id.String,
			ProcutID:        ProcutID.String,
			SaleID:          SaleID.String,
			SaleIncrementID: SaleIncrementID.String,
			Quantity:        int(Quantity.Int64),
			Price:           Price.Float64,
			TotalPrice:      TotalPrice.Float64,
			CreatedAt:       CreatedAt.String,
			UpdatedAt:       UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (r *saleProductRepo) Update(ctx context.Context, req *models.UpdateSaleProduct) (int64, error) {

	query := `
		UPDATE "sale"
			SET
				"product_id" = $2,
				"sale_id" = $3,
				"sale_increment_id" = $4,
				"quantity" = $5,
				"price" = $6,
				"total_price" = $7,
				"updated_at" = NOW()
		WHERE "id" = $1
	`
	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.ProcutID,
		req.SaleID,
		req.SaleIncrementID,
		req.Quantity,
		req.Price,
		float64(req.Quantity)*req.Price,
		req.TotalPrice,
	)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}

func (r *saleProductRepo) Delete(ctx context.Context, req *models.SaleProductPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sale_product WHERE id = $1", req.Id)
	return err
}
