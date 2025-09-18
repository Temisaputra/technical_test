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

type MaterialRepository struct {
	*TransactionRepository
}

func NewMaterialRepo(db *gorm.DB) irepository.MaterialRepository {
	return &MaterialRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *MaterialRepository) GetAllMaterial(ctx context.Context, pagination *request.Pagination, params *presenter.MaterialRequest) (materials []*presenter.MaterialResponse, meta response.Meta, err error) {
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.MaterialList{})

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

	return materials, meta, nil
}

func (r *MaterialRepository) GetMaterialByID(ctx context.Context, id int) (material *presenter.MaterialResponse, err error) {
	var result entity.MaterialList
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.MaterialList{}).Where("id = ?", id)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
		return nil, fmt.Errorf("[MaterialRepository-GetMaterialByID] Error when getting material by id: %w", err)
	}

	return result.ToPresenter(), nil
}

func (r *MaterialRepository) CreateMaterial(ctx context.Context, params *[]entity.MaterialList) error {
	conn := r.Conn(ctx)
	currentTime := time.Now()

	for i := range *params {
		(*params)[i].CreatedAt = currentTime
	}

	db := conn.WithContext(ctx).Model(&entity.MaterialList{})
	if err := db.Create(params).Error; err != nil {
		return fmt.Errorf("[MaterialRepository-CreateMaterial] Error when creating material: %w", err)
	}
	return nil
}

func (r *MaterialRepository) UpdateMaterial(ctx context.Context, params *[]entity.MaterialList, id int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now()

	for i := range *params {
		(*params)[i].UpdatedAt = &currentTime
		(*params)[i].SupplierID = uint(id)
	}

	db := conn.WithContext(ctx).Model(&entity.MaterialList{})
	if err := db.Save(params).Error; err != nil {
		return fmt.Errorf("[MaterialRepository-UpdateMaterial] Error when updating material: %w", err)
	}
	return nil

}

func (r *MaterialRepository) DeleteMaterial(ctx context.Context, id int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.MaterialList{})
	if err := db.Where("id = ?", id).Delete(&entity.MaterialList{}).Error; err != nil {
		return err
	}

	return nil
}
func (r *MaterialRepository) GetMaterialByIdSupplier(ctx context.Context, supplierID int) (*[]entity.MaterialList, error) {
	var result []entity.MaterialList
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.MaterialList{}).Where("supplier_id = ?", supplierID)
	if err := db.Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[MaterialRepository-GetMaterialByIdSupplier] Material not found: %w", err)
		}
		return nil, fmt.Errorf("[MaterialRepository-GetMaterialByIdSupplier] Error when getting material by supplier id: %w", err)
	}

	return &result, nil
}

func (r *MaterialRepository) DeleteMaterialByIdSupplier(ctx context.Context, supplierID int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.MaterialList{})
	if err := db.Where("supplier_id = ?", supplierID).Delete(&entity.MaterialList{}).Error; err != nil {
		return err
	}

	return nil
}
