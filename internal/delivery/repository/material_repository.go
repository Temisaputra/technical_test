package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type MaterialRepository interface {
	GetAllMaterial(ctx context.Context, pagination *request.Pagination, params *presenter.MaterialRequest) ([]*presenter.MaterialResponse, response.Meta, error)
	GetMaterialByID(ctx context.Context, id int) (*presenter.MaterialResponse, error)
	CreateMaterial(ctx context.Context, material *[]entity.MaterialList) error
	UpdateMaterial(ctx context.Context, material *[]entity.MaterialList, id int) error
	DeleteMaterial(ctx context.Context, id int) error
	GetMaterialByIdSupplier(ctx context.Context, supplierID int) (*[]entity.MaterialList, error)
	DeleteMaterialByIdSupplier(ctx context.Context, supplierID int) error
}
