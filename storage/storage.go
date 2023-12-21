package storage

import (
	"context"
	"time"

	"exam/models"
)

type StorageI interface {
	Coming() ComingRepoI
	Branch() BranchRepoI
	Client() ClientRepoI
	Product() ProductRepoI
	SaleProduct() SaleProductRepoI
	Remainder() RemainderRepoI
	Sale() SaleRepoI
	PickingList() PickingListRepoI
	IncrementID() IncrementIDRepoI
	Report() ReportRepoI
}

type ComingRepoI interface {
	Create(ctx context.Context, req *models.CreateComing) (*models.Coming, error)
	GetByID(ctx context.Context, req *models.ComingPrimaryKey) (*models.Coming, error)
	GetList(ctx context.Context, req *models.GetListComingRequest) (*models.GetListComingResponse, error)
	Update(ctx context.Context, req *models.UpdateComing) (int64, error)
	Delete(ctx context.Context, req *models.ComingPrimaryKey) error
}

type ProductRepoI interface {
	Create(ctx context.Context, req *models.CreateProduct) (*models.Product, error)
	GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error)
	GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error)
	Update(ctx context.Context, req *models.UpdateProduct) (int64, error)
	Delete(ctx context.Context, req *models.ProductPrimaryKey) error
}

type SaleRepoI interface {
	Create(ctx context.Context, req *models.CreateSale) (*models.Sale, error)
	GetByID(ctx context.Context, req *models.SalePrimaryKey) (*models.Sale, error)
	GetList(ctx context.Context, req *models.GetListSaleRequest) (*models.GetListSaleResponse, error)
	Update(ctx context.Context, req *models.UpdateSale) (int64, error)
	Delete(ctx context.Context, req *models.SalePrimaryKey) error
}

type SaleProductRepoI interface {
	Create(ctx context.Context, req *models.CreateSaleProduct) (*models.SaleProduct, error)
	GetByID(ctx context.Context, req *models.SaleProductPrimaryKey) (*models.SaleProduct, error)
	GetList(ctx context.Context, req *models.GetListSaleProductRequest) (*models.GetListSaleProductResponse, error)
	Update(ctx context.Context, req *models.UpdateSaleProduct) (int64, error)
	Delete(ctx context.Context, req *models.SaleProductPrimaryKey) error
}

type RemainderRepoI interface {
	Create(ctx context.Context, req *models.Remainder) (*models.Remainder, error)
	GetByID(ctx context.Context, req *models.RemainderPrimaryKey) (*models.Remainder, error)
	GetList(ctx context.Context, req *models.GetListRemainderRequest) (*models.GetListRemainderResponse, error)
	Update(ctx context.Context, req *models.Remainder) (int64, error)
	Delete(ctx context.Context, req *models.RemainderPrimaryKey) error
}

type BranchRepoI interface {
	Create(ctx context.Context, req *models.CreateBranch) (*models.Branch, error)
	GetByID(ctx context.Context, req *models.BranchPrimaryKey) (*models.Branch, error)
	GetList(ctx context.Context, req *models.GetListBranchRequest) (*models.GetListBranchResponse, error)
	Update(ctx context.Context, req *models.UpdateBranch) (int64, error)
	Delete(ctx context.Context, req *models.BranchPrimaryKey) error
}

type ClientRepoI interface {
	Create(ctx context.Context, req *models.CreateClient) (*models.Client, error)
	GetByID(ctx context.Context, req *models.ClientPrimaryKey) (*models.Client, error)
	GetList(ctx context.Context, req *models.GetListClientRequest) (*models.GetListClientResponse, error)
	Update(ctx context.Context, req *models.UpdateClient) (int64, error)
	Delete(ctx context.Context, req *models.ClientPrimaryKey) error
}

type PickingListRepoI interface {
	Create(ctx context.Context, req *models.PickingList) (*models.PickingList, error)
	GetByID(ctx context.Context, req *models.PickingListPrimaryKey) (*models.PickingList, error)
	GetList(ctx context.Context, req *models.GetListPickingListRequest) (*models.GetListPickingListResponse, error)
	Update(ctx context.Context, req *models.PickingList) (int64, error)
	Delete(ctx context.Context, req *models.PickingListPrimaryKey) error
}
type IncrementIDRepoI interface {
	GetLast(ctx context.Context, tableName string, columnName string) (string, error)
}

type ReportRepoI interface {
	GetListReport(context.Context, time.Time, time.Time) (*models.GetListClientResponse, error)
	GetListSaleBranch(context.Context, time.Time, time.Time) (*models.AllSaleReport, error)
}
