package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type SupplierRepository interface {
	GetAllSupplier(ctx context.Context, pagination *request.Pagination, params *presenter.SupplierRequest) ([]*presenter.SupplierResponse, response.Meta, error)
	GetSupplierByID(ctx context.Context, id int) (*presenter.SupplierResponse, error)
	CreateSupplier(ctx context.Context, supplier *presenter.SupplierRequest) (int, error)
	UpdateSupplier(ctx context.Context, supplier *presenter.SupplierRequest, id int) error
	DeleteSupplier(ctx context.Context, id int) error
	BlockUnblockSupplier(ctx context.Context, id int, status string) error
	InsertHistorySupplier(ctx context.Context, history *entity.StatusHistory) error
	ExportSupplier(ctx context.Context, params *presenter.SupplierRequest) (exportData []presenter.SupplierResponse, err error)
}
