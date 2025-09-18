package repository

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type ApprovalRepository interface {
	GetAllApprovalWorkflows(ctx context.Context, pagination *request.Pagination, params *presenter.ApprovalRequest) ([]*presenter.ApprovalResponse, response.Meta, error)
	GetActiveWorkflowBySupplierID(ctx context.Context, supplierID uint) (*entity.ApprovalWorkflow, error)
	CreateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) error
	UpdateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) error
	GetLogsByWorkflowID(ctx context.Context, workflowID uint) ([]entity.ApprovalLog, error)
	CreateLog(ctx context.Context, log *entity.ApprovalLog) error
}
