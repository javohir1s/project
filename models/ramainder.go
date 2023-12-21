package models

type RemainderPrimaryKey struct {
	Id string `json:"id"`
}

type CreateRemainder struct {
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	ComingPrice float64 `json:"coming_price"`
	SalePrice   float64 `json:"sale_price"`
	BranchID    string  `json:"branch_id"`
}

type Remainder struct {
	Id          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	ComingPrice float64 `json:"coming_price"`
	SalePrice   float64 `json:"sale_price"`
	BranchID    string  `json:"branch_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateRemainder struct {
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	ComingPrice float64 `json:"coming_price"`
	SalePrice   float64 `json:"sale_price"`
	BranchID    string  `json:"branch_id"`
}

type GetListRemainderRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListRemainderResponse struct {
	Count      int          `json:"count"`
	Remainders []*Remainder `json:"remainders"`
}