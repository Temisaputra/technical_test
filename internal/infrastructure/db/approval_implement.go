package db

import (
	"context"
	"strings"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	irepository "github.com/Temisaputra/warOnk/internal/delivery/repository"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	*TransactionRepository
}

func NewApprovalRepo(db *gorm.DB) irepository.ApprovalRepository {
	return &ApprovalRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *ApprovalRepository) GetAllApprovalWorkflows(
	ctx context.Context,
	pagination *request.Pagination,
	params *presenter.ApprovalRequest,
) (workflows []*presenter.ApprovalResponse, meta response.Meta, err error) {

	db := r.Conn(ctx).WithContext(ctx).Model(&entity.ApprovalWorkflow{})

	// --- Filtering dari params ---
	if params != nil {
		if params.Status != "" {
			db = db.Where("status = ?", params.Status)
		}
		if params.Stage != "" {
			db = db.Where("stage = ?", params.Stage)
		}
		if params.SupplierID > 0 {
			db = db.Where("supplier_id = ?", params.SupplierID)
		}
	}

	// --- Sorting ---
	if pagination.OrderBy != "" {
		order := pagination.OrderBy
		if pagination.OrderType != "" {
			order += " " + pagination.OrderType
		}
		// jika ada format "field|desc" → override
		if strings.Contains(pagination.OrderBy, "|") {
			parts := strings.Split(pagination.OrderBy, "|")
			if len(parts) == 2 {
				order = parts[0] + " " + parts[1]
			}
		}
		db = db.Order(order)
	} else {
		db = db.Order("updated_at DESC NULLS LAST")
	}

	// --- Count total data ---
	if err = db.Count(&meta.TotalData).Error; err != nil {
		return nil, meta, err
	}

	// --- Pagination ---
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	var entities []entity.ApprovalWorkflow
	if err = db.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, meta, err
	}

	// --- Mapping entity ke response ---
	for _, e := range entities {
		workflows = append(workflows, &presenter.ApprovalResponse{
			ID:         e.ID,
			SupplierID: e.SupplierID,
			Stage:      e.Stage,
			SLAHours:   e.SLAHours,
			Status:     e.Status,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		})
	}

	// --- Meta info ---
	meta.Page = int64(pagination.Page)
	meta.PageSize = int64(pagination.PageSize)
	if meta.TotalData > 0 {
		meta.TotalPage = (meta.TotalData + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)
	}

	return workflows, meta, nil
}

func (r *ApprovalRepository) GetActiveWorkflowBySupplierID(ctx context.Context, supplierID uint) (*entity.ApprovalWorkflow, error) {
	var workflow entity.ApprovalWorkflow
	err := r.Conn(ctx).WithContext(ctx).
		Where("supplier_id = ? AND status = ?", supplierID, "In Progress").
		Order("created_at DESC").
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *ApprovalRepository) CreateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) error {
	return r.Conn(ctx).WithContext(ctx).Create(workflow).Error
}

func (r *ApprovalRepository) UpdateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) error {
	return r.Conn(ctx).WithContext(ctx).Save(workflow).Error
}

func (r *ApprovalRepository) GetLogsByWorkflowID(ctx context.Context, workflowID uint) ([]entity.ApprovalLog, error) {
	var logs []entity.ApprovalLog
	err := r.Conn(ctx).WithContext(ctx).
		Where("workflow_id = ?", workflowID).
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}

func (r *ApprovalRepository) CreateLog(ctx context.Context, log *entity.ApprovalLog) error {
	return r.Conn(ctx).WithContext(ctx).Create(log).Error
}
