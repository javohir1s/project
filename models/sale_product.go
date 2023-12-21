package models

type SaleProductPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSaleProduct struct {
	ProcutID        string  `json:"product_id"`
	SaleID          string  `json:"sale_id"`
	SaleIncrementID string  `json:"sale_increment_id"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	TotalPrice      float64 `json:"total_price"`
}

type SaleProduct struct {
	Id              string  `json:"id"`
	ProcutID        string  `json:"product_id"`
	SaleID          string  `json:"sale_id"`
	SaleIncrementID string  `json:"sale_increment_id"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	TotalPrice      float64 `json:"total_price"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type UpdateSaleProduct struct {
	Id              string  `json:"id"`
	ProcutID        string  `json:"product_id"`
	SaleID          string  `json:"sale_id"`
	SaleIncrementID string  `json:"sale_increment_id"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	TotalPrice      float64 `json:"total_price"`
}

type GetListSaleProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSaleProductResponse struct {
	Count        int            `json:"count"`
	SaleProducts []*SaleProduct `json:"sale_products"`
}
