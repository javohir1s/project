package models

type SalePrimaryKey struct {
	Id string `json:"id"`
}

type CreateSale struct {
	ClientID    string  `json:"client_id"`
	BranchID    string  `json:"branch_id"`
	IncrementID string  `json:"increment_id"`
	TotalPrice  float64 `json:"total_price"`
}

type Sale struct {
	Id          string  `json:"id"`
	ClientID    string  `json:"client_id"`
	BranchID    string  `json:"branch_id"`
	IncrementID string  `json:"increment_id"`
	TotalPrice  float64 `json:"total_price"`
	Paid        float64 `json:"paid"`
	Debd        float64 `json:"debd"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateSale struct {
	Id          string  `json:"id"`
	ClientID    string  `json:"client_id"`
	BranchID    string  `json:"branch_id"`
	IncrementID string  `json:"increment_id"`
	TotalPrice  float64 `json:"total_price"`
	Paid        float64 `json:"paid"`
	Debd        float64 `json:"debd"`
}

type GetListSaleRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListSaleResponse struct {
	Count int     `json:"count"`
	Sales []*Sale `json:"sales"`
}
