package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type GroupRepository interface {
	GetAllGroups(ctx context.Context, pagination *request.Pagination) ([]*presenter.GroupResponse, response.Meta, error)
	GetGroupByID(ctx context.Context, id int) (*presenter.GroupResponse, error)
	CreateGroup(ctx context.Context, group *[]entity.SupplierGroup) error
	UpdateGroup(ctx context.Context, group *entity.SupplierGroup, id int) error
	DeleteGroup(ctx context.Context, id int) error
	GetGroupByIdSupplier(ctx context.Context, supplierID int) (*entity.Group, error)
	DeleteGroupByIdSupplier(ctx context.Context, supplierID int) error
	CreateGroups(ctx context.Context, groups presenter.GroupRequest) error
	UpdateGroups(ctx context.Context, groups presenter.GroupRequest, groupID int) error
}
