package models

type PickingList struct {
	ID                string  `json:"id"`
	Product_ID        string  `json:"product_id"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	Total_price       float64 `json:"total_price"`
	ComingID          string  `json:"coming_id"`
	ComingIncrementID string  `json:"coming_increment_id"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type CreatePickingList struct {
	Product_ID        string  `json:"product_id"`
	Quantity          int     `json:"quantity"`
	Price             float64 `json:"price"`
	ComingIncrementID string  `json:"coming_increment_id"`
}

type PickingListPrimaryKey struct {
	Id string `json:"id"`
}

type UpdatePickingList struct {
	Product_ID        string  `json:"product_id"`
	ComingID          string  `json:"coming_id"`
	Price             float64 `json:"price"`
	ComingIncrementID string  `json:"coming_increment_id"`
	Quantity          int     `json:"quantity"`
}

type GetListPickingListRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListPickingListResponse struct {
	Count     int            `json:"count"`
	Pickinges []*PickingList `json:"picking_list"`
}