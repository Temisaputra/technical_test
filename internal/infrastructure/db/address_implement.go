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

type AddressRepository struct {
	*TransactionRepository
}

func NewAddressRepo(db *gorm.DB) irepository.AddressRepository {
	return &AddressRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *AddressRepository) GetAllAddress(ctx context.Context, pagination *request.Pagination, params *presenter.AddressRequest) (addresses []*presenter.AddressResponse, meta response.Meta, err error) {
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

	return addresses, meta, nil
}

func (r *AddressRepository) GetAddressByID(ctx context.Context, id int) (address *presenter.AddressResponse, err error) {
	var result entity.Address
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Address{}).Where("id = ?", id)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[AddressRepository-GetAddressByID] Address not found: %w", err)
		}
		return nil, fmt.Errorf("[AddressRepository-GetAddressByID] Error when getting address by id: %w", err)
	}

	return result.ToPresenter(), nil
}

func (r *AddressRepository) CreateAddress(ctx context.Context, params *entity.Address) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	storedData := map[string]interface{}{
		"supplier_id": params.SupplierID,
		"name":        params.Name,
		"address":     params.Address,
		"main":        params.Main,
		"created_at":  currentTime,
	}
	db := conn.WithContext(ctx).Model(&entity.Address{})
	if err := db.Create(&storedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) UpdateAddress(ctx context.Context, params *entity.Address, id int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	updatedData := map[string]interface{}{
		"name":       params.Name,
		"address":    params.Address,
		"main":       params.Main,
		"updated_at": currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Address{})
	if err := db.Where("id = ?", id).Updates(&updatedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) DeleteAddress(ctx context.Context, id int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.Address{})
	if err := db.Where("id = ?", id).Delete(&entity.Address{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) GetAddressByIdSupplier(ctx context.Context, supplierID int) (*entity.Address, error) {
	var result entity.Address
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Address{}).Where("supplier_id = ?", supplierID).Where("main = ?", true)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[AddressRepository-GetAddressByIdSupplier] Address not found: %w", err)
		}
		return nil, fmt.Errorf("[AddressRepository-GetAddressByIdSupplier] Error when getting address by supplier id: %w", err)
	}

	return &result, nil
}

func (r *AddressRepository) DeleteAddressByIdSupplier(ctx context.Context, supplierID int) error {
	conn := r.Conn(ctx).WithContext(ctx)
	db := conn.WithContext(ctx).Model(&entity.Address{})
	if err := db.Where("supplier_id = ?", supplierID).Delete(&entity.Address{}).Error; err != nil {
		return err
	}
	return nil
}
