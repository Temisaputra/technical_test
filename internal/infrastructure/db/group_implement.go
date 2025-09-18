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

type GroupRepository struct {
	*TransactionRepository
}

func NewGroupRepo(db *gorm.DB) irepository.GroupRepository {
	return &GroupRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *GroupRepository) GetAllGroups(ctx context.Context, pagination *request.Pagination) (groups []*presenter.GroupResponse, meta response.Meta, err error) {
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Group{})

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

	var result []entity.Group

	if err = db.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, meta, err
	}

	meta.Page = int64(pagination.Page)
	meta.PageSize = int64(pagination.PageSize)
	meta.TotalPage = meta.TotalData/int64(pagination.PageSize) + 1

	for _, v := range result {
		groups = append(groups, v.ToPresenter())
	}

	return groups, meta, nil
}

func (r *GroupRepository) GetGroupByID(ctx context.Context, id int) (group *presenter.GroupResponse, err error) {
	var result entity.Group
	db := r.Conn(ctx).WithContext(ctx).Model(&entity.Group{}).Where("id = ?", id)
	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[GroupRepository-GetGroupByID] Group not found: %w", err)
		}
		return nil, fmt.Errorf("[GroupRepository-GetGroupByID] Error when getting group by id: %w", err)
	}

	return result.ToPresenter(), nil
}

func (r *GroupRepository) CreateGroup(ctx context.Context, params *[]entity.SupplierGroup) error {
	conn := r.Conn(ctx)
	currentTime := time.Now()
	for i := range *params {
		(*params)[i].CreatedAt = currentTime
	}

	db := conn.WithContext(ctx).Model(&entity.SupplierGroup{})
	if err := db.Create(params).Error; err != nil {
		return fmt.Errorf("[GroupRepository-CreateGroup] Error when creating group: %w", err)
	}
	return nil
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, params *entity.SupplierGroup, id int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now()

	updatedData := map[string]interface{}{
		"supplier_id": params.SupplierID,
		"group_id":    params.GroupID,
		"updated_at":  currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Group{})
	if err := db.Where("id = ?", id).Updates(&updatedData).Error; err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) DeleteGroup(ctx context.Context, id int) error {
	conn := r.Conn(ctx).WithContext(ctx)

	db := conn.WithContext(ctx).Model(&entity.Group{})
	if err := db.Where("id = ?", id).Delete(&entity.Group{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) GetGroupByIdSupplier(ctx context.Context, supplierID int) (*entity.Group, error) {
	var result entity.Group
	db := r.Conn(ctx).WithContext(ctx).
		Table("supplier_groups").
		Select("groups.*").
		Joins("JOIN groups ON groups.id = supplier_groups.group_id").
		Where("supplier_groups.supplier_id = ? AND groups.is_active = ?", supplierID, true).
		Find(&result)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[GroupRepository-GetGroupByIdSupplier] Group not found: %w", err)
		}
		return nil, fmt.Errorf("[GroupRepository-GetGroupByIdSupplier] Error when getting group by supplier id: %w", err)
	}

	return &result, nil
}

func (r *GroupRepository) DeleteGroupByIdSupplier(ctx context.Context, supplierID int) error {
	conn := r.Conn(ctx).WithContext(ctx)
	db := conn.WithContext(ctx).Model(&entity.SupplierGroup{})
	if err := db.Where("supplier_id = ?", supplierID).Delete(&entity.SupplierGroup{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GroupRepository) CreateGroups(ctx context.Context, groups presenter.GroupRequest) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	storedData := map[string]interface{}{
		"name":       groups.Name,
		"value":      groups.Description,
		"is_active":  groups.Active,
		"created_at": currentTime,
		"updated_at": currentTime,
	}

	db := conn.WithContext(ctx).Model(&entity.Group{})
	if err := db.Create(&storedData).Error; err != nil {
		return fmt.Errorf("[GroupRepository-CreateGroups] Error when creating groups: %w", err)
	}

	return nil
}

func (r *GroupRepository) UpdateGroups(ctx context.Context, groups presenter.GroupRequest, groupID int) error {
	conn := r.Conn(ctx)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	updatedData := map[string]interface{}{
		"name":       groups.Name,
		"value":      groups.Description,
		"is_active":  groups.Active,
		"updated_at": currentTime,
	}
	db := conn.WithContext(ctx).Model(&entity.Group{})
	if err := db.Where("id = ?", groupID).Updates(&updatedData).Error; err != nil {
		return fmt.Errorf("[GroupRepository-UpdateGroups] Error when updating groups: %w", err)
	}
	return nil
}
