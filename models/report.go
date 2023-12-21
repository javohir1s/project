package models

type OverallReport struct {
	Clients           []*Client     `json:"clients_report"`
	BranchSaleReports AllSaleReport `json:"sale_branches_report"`
}
type AllSaleReport struct {
	SaleReports []*SaleReport
}
type SaleReport struct {
	Name     string  `json:"name"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}
