package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	irepository "github.com/Temisaputra/warOnk/internal/delivery/repository"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"gorm.io/gorm"
)

type ContactRepository struct {
	*TransactionRepository
}

func NewContactRepo(db *gorm.DB) irepository.ContactRepository {
	return &ContactRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *ContactRepository) GetAllContact(ctx context.Context, pagination *request.Pagination, params *presenter.ContactRequest) (contacts []*presenter.ContactResponse, meta response.Meta, err error) {
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Address{})

	if pagination.OrderBy != "" || pagination.OrderType != "" {
		splitOrder := strings.Split(pagination.OrderBy, "|")
		if len(splitOrder) > 1 {
			db = db.Order(splitOrder[0] + " " + splitOrder[1] + " " + pagination.OrderType)
		} else {
			db = db.Order(pagination.OrderBy + " " + pagination.OrderType)
		}
	} else {
		db = db.Order("updated_at DESC NULLS LAST")
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	if err = db.Count(&meta.TotalData).Error; err != nil {
		return nil, meta, err
	}

	var result []entity.Address

	if err = db.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, meta, err
	}

	meta.Page = int64(pagination.Page)
	meta.PageSize = int64(pagination.PageSize)
	meta.TotalPage = meta.TotalData/int64(pagination.PageSize) + 1

	return contacts, meta, nil
}

func (r *ContactRepository) GetContactByID(ctx context.Context, id int) (contact *presenter.ContactResponse, err error) {
	var result entity.Contact
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Contact{}).Where("id = ?", id)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[ContactRepository-GetContactByID] Contact not found: %w", err)
		}
		return nil, fmt.Errorf("[ContactRepository-GetContactByID] Error when getting contact by id: %w", err)
	}

	return result.ToPresenter(), nil
}

func (r *ContactRepository) CreateContact(ctx context.Context, params *entity.Contact) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	storedData := map[string]interface{}{
		"supplier_id":  params.SupplierID,
		"name":         params.Name,
		"job_position": params.JobPosition,
		"email":        params.Email,
		"phone":        params.Phone,
		"mobile":       params.Mobile,
		"main":         params.Main,
		"created_at":   currentTime,
	}
	db := conn.WithContext(ctx).Model(&entity.Contact{})
	if err := db.Create(&storedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *ContactRepository) UpdateContact(ctx context.Context, params *entity.Contact, id int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	updatedData := map[string]interface{}{
		"name":         params.Name,
		"job_position": params.JobPosition,
		"email":        params.Email,
		"phone":        params.Phone,
		"mobile":       params.Mobile,
		"main":         params.Main,
		"updated_at":   currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Contact{})
	if err := db.Where("id = ?", id).Updates(&updatedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *ContactRepository) DeleteContact(ctx context.Context, id int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.Contact{})
	if err := db.Where("id = ?", id).Delete(&entity.Contact{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *ContactRepository) GetContactByIdSupplier(ctx context.Context, supplierID int) (*entity.Contact, error) {
	var result entity.Contact
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Contact{}).Where("supplier_id = ?", supplierID).Where("main = ?", true)
	if err := db.First(&result).Error; err != nil {
		return nil, fmt.Errorf("[ContactRepository-GetContactByIdSupplier] Error when getting contact by supplier id: %w", err)
	}
	return &result, nil
}

func (r *ContactRepository) DeleteContactByIdSupplier(ctx context.Context, supplierID int) error {
	conn := r.Conn(ctx).WithContext(ctx)
	db := conn.WithContext(ctx).Model(&entity.Contact{})
	if err := db.Where("supplier_id = ?", supplierID).Delete(&entity.Contact{}).Error; err != nil {
		return err
	}
	return nil
}
