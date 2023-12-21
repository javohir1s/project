package models

type ComingPrimaryKey struct {
	Id string `json:"id"`
}

type CreateComing struct {
	IncrementID string `json:"increment_id"`
	Date        string `json:"date"`
	BranchID    string `json:"branch_id"`
}

type Coming struct {
	Id          string `json:"id"`
	Date        string `json:"date"`
	IncrementID string `json:"increment_id"`
	BranchID    string `json:"branch_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateComing struct {
	Id       string `json:"id"`
	Date     string `json:"date"`
	BranchID string `json:"branch_id"`
	Status   string `json:"status"`
}

type GetListComingRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListComingResponse struct {
	Count    int       `json:"count"`
	Cominges []*Coming `json:"cominges"`
}
