package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type ContactRepository interface {
	GetAllContact(ctx context.Context, pagination *request.Pagination, params *presenter.ContactRequest) ([]*presenter.ContactResponse, response.Meta, error)
	GetContactByID(ctx context.Context, id int) (*presenter.ContactResponse, error)
	CreateContact(ctx context.Context, contact *entity.Contact) error
	UpdateContact(ctx context.Context, contact *entity.Contact, id int) error
	DeleteContact(ctx context.Context, id int) error
	GetContactByIdSupplier(ctx context.Context, supplierID int) (*entity.Contact, error)
	DeleteContactByIdSupplier(ctx context.Context, supplierID int) error
}
