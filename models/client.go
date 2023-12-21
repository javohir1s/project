package models

type ClientPrimaryKey struct {
	Id string `json:"id"`
}

type CreateClient struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FatherName string `json:"father_name"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	Gender     string `json:"gender"`
	BranchID   string `json:"branch_id"`
	Active     bool   `json:"active"`
}

type Client struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FatherName string `json:"father_name"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	Gender     string `json:"gender"`
	BranchID   string `json:"branch_id"`
	Active     bool   `json:"active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UpdateClient struct {
	Id         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FatherName string `json:"father_name"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	Gender     string `json:"gender"`
	BranchID   string `json:"branch_id"`
	Active     bool   `json:"active"`
}

type GetListClientRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
	Query  string `json:"query"`
}

type GetListClientResponse struct {
	Count   int       `json:"count"`
	Clients []*Client `json:"clients"`
}
