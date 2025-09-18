package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type AddressRepository interface {
	GetAllAddress(ctx context.Context, pagination *request.Pagination, params *presenter.AddressRequest) ([]*presenter.AddressResponse, response.Meta, error)
	GetAddressByID(ctx context.Context, id int) (*presenter.AddressResponse, error)
	CreateAddress(ctx context.Context, address *entity.Address) error
	UpdateAddress(ctx context.Context, address *entity.Address, id int) error
	DeleteAddress(ctx context.Context, id int) error
	GetAddressByIdSupplier(ctx context.Context, supplierID int) (*entity.Address, error)
	DeleteAddressByIdSupplier(ctx context.Context, supplierID int) error
}
