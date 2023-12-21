package postgres

import (
	"context"
	"fmt"

	"exam/config"
	"exam/pkg/helpers"
	"exam/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db          *pgxpool.Pool
	coming      storage.ComingRepoI
	branch      storage.BranchRepoI
	client      storage.ClientRepoI
	product     storage.ProductRepoI
	sale        storage.SaleRepoI
	saleProduct storage.SaleProductRepoI
	remainder   storage.RemainderRepoI
	pickingList storage.PickingListRepoI
	getIncrementId storage.IncrementIDRepoI
	report        storage.ReportRepoI

}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	config, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
			cfg.PostgresHost,
			cfg.PostgresUser,
			cfg.PostgresDatabase,
			cfg.PostgresPassword,
			cfg.PostgresPort,
		),
	)

	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnection

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		return nil, err

	}

	return &Store{
		db: pgxpool,
	}, nil
}

func (s *Store) Remainder() storage.RemainderRepoI {

	if s.remainder == nil {
		s.remainder = NewRemainderRepo(s.db)
	}

	return s.remainder
}
func (s *Store) Client() storage.ClientRepoI {

	if s.client == nil {
		s.client = NewClientRepo(s.db)
	}

	return s.client
}
func (s *Store) Product() storage.ProductRepoI {

	if s.product == nil {
		s.product = NewProductRepo(s.db)
	}

	return s.product
}
func (s *Store) PickingList() storage.PickingListRepoI {

	if s.pickingList == nil {
		s.pickingList = NewPickingListRepo(s.db)
	}

	return s.pickingList
}

func (s *Store) Sale() storage.SaleRepoI {

	if s.sale == nil {
		s.sale = NewSaleRepo(s.db)
	}

	return s.sale
}
func (s *Store) SaleProduct() storage.SaleProductRepoI {

	if s.saleProduct == nil {
		s.saleProduct = NewSaleProductRepo(s.db)
	}

	return s.saleProduct
}
func (s *Store) Coming() storage.ComingRepoI {

	if s.coming == nil {
		s.coming = NewComingRepo(s.db)
	}

	return s.coming
}

func (s *Store) Branch() storage.BranchRepoI {

	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}

	return s.branch
}

func (s *Store) IncrementID() storage.IncrementIDRepoI {

	if s.getIncrementId == nil {
		s.getIncrementId = helpers.NewIncrementRepo(s.db)
	}

	return s.getIncrementId
}

func (s *Store) Report() storage.ReportRepoI {
	if s.report == nil {
		s.report = NewReportRepo(s.db)
	}
	return s.report
}
