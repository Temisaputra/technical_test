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

type SupplierRepository struct {
	*TransactionRepository
}

func NewSupplierRepo(db *gorm.DB) irepository.SupplierRepository {
	return &SupplierRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *SupplierRepository) GetAllSupplier(ctx context.Context, pagination *request.Pagination, params *presenter.SupplierRequest) (suppliers []*presenter.SupplierResponse, meta response.Meta, err error) {
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Supplier{}).
		Joins("JOIN addresses ON addresses.supplier_id = suppliers.id").
		Joins("JOIN contacts ON contacts.supplier_id = suppliers.id").
		Group("suppliers.id, addresses.supplier_id, contacts.supplier_id")

	if pagination.Keyword != "" {
		keywordStr := "%" + pagination.Keyword + "%"
		db = db.Where("id ILIKE ? OR name ILIKE ? OR nick_name ILIKE ? OR status ILIKE ? ", keywordStr, keywordStr, keywordStr, keywordStr)
	}

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	if pagination.OrderBy != "" || pagination.OrderType != "" {
		splitOrder := strings.Split(pagination.OrderBy, "|")
		if len(splitOrder) > 1 {
			db = db.Order(splitOrder[0] + " " + splitOrder[1] + " " + pagination.OrderType)
		} else {
			db = db.Order(pagination.OrderBy + " " + pagination.OrderType)
		}
	} else {
		db = db.Order("suppliers.updated_at DESC NULLS LAST")
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	if err = db.Count(&meta.TotalData).Error; err != nil {
		return nil, meta, err
	}

	var result []entity.Supplier

	if err = db.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, meta, err
	}

	meta.Page = int64(pagination.Page)
	meta.PageSize = int64(pagination.PageSize)
	meta.TotalPage = meta.TotalData/int64(pagination.PageSize) + 1

	for _, item := range result {
		suppliers = append(suppliers, item.ToPresenter())
	}

	return
}

func (r *SupplierRepository) GetSupplierByID(ctx context.Context, id int) (supplier *presenter.SupplierResponse, err error) {
	var result entity.Supplier
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Supplier{}).Where("id = ?", id)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[SupplierRepository-GetSupplierByID] Supplier not found: %w", err)
		}
		return nil, fmt.Errorf("[SupplierRepository-GetSupplierByID] Error when getting supplier by id: %w", err)
	}

	return result.ToPresenter(), nil
}

func (r *SupplierRepository) CreateSupplier(ctx context.Context, params *presenter.SupplierRequest) (int, error) {
	conn := r.Conn(ctx)
	currentTime := time.Now()

	supplier := entity.Supplier{
		Name:      params.Name,
		NickName:  params.NickName,
		Status:    params.Status,
		CreatedAt: currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Supplier{})
	if err := db.Create(&supplier).Error; err != nil {
		return 0, err
	}

	return int(supplier.ID), nil
}

func (r *SupplierRepository) UpdateSupplier(ctx context.Context, params *presenter.SupplierRequest, id int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	updatedData := map[string]interface{}{
		"name":       params.Name,
		"nick_name":  params.NickName,
		"status":     params.Status,
		"updated_at": currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Supplier{})
	if err := db.Where("id = ?", id).Updates(&updatedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupplierRepository) DeleteSupplier(ctx context.Context, id int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.Supplier{})
	if err := db.Where("id = ?", id).Delete(&entity.Supplier{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *SupplierRepository) BlockUnblockSupplier(ctx context.Context, id int, status string) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var newStatus string
	if status == "Blocked" {
		newStatus = "Active"
	} else {
		newStatus = "Blocked"
	}

	updatedData := map[string]interface{}{
		"updated_at": currentTime,
		"status":     newStatus,
	}

	db := conn.WithContext(ctx).Model(&entity.Supplier{})
	if err := db.Where("id = ?", id).Updates(&updatedData).Error; err != nil {
		return err
	}
	return nil
}

func (r *SupplierRepository) InsertHistorySupplier(ctx context.Context, history *entity.StatusHistory) error {
	conn := r.Conn(ctx)
	currentTime := time.Now()
	history.CreatedAt = currentTime

	db := conn.WithContext(ctx).Model(&entity.StatusHistory{})
	if err := db.Create(history).Error; err != nil {
		return fmt.Errorf("[SupplierRepository-InsertHistorySupplier] Error when inserting history supplier: %w", err)
	}

	return nil
}

func (r *SupplierRepository) ExportSupplier(ctx context.Context, params *presenter.SupplierRequest) (suppliers []presenter.SupplierResponse, err error) {
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Supplier{}).
		Joins("JOIN addresses ON addresses.supplier_id = suppliers.id").
		Joins("JOIN contacts ON contacts.supplier_id = suppliers.id").
		Group("suppliers.id, addresses.supplier_id, contacts.supplier_id")

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	var result []entity.Supplier

	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}

	for _, item := range result {
		suppliers = append(suppliers, *item.ToPresenter())
	}

	return
}
