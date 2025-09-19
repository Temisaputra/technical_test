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
	GetActiveWorkflowByID(ctx context.Context, Id uint) (*entity.ApprovalWorkflow, error)
	CreateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) (uint, error)
	UpdateWorkflow(ctx context.Context, workflow *entity.ApprovalWorkflow) error
	GetLogsByWorkflowID(ctx context.Context, workflowID uint) ([]entity.ApprovalLog, error)
	CreateLog(ctx context.Context, log *entity.ApprovalLog) error
	CreateWorkflowSteps(ctx context.Context, steps []entity.ApprovalStep) error
	UpdateWorkflowStep(ctx context.Context, step *entity.ApprovalStep) error
	GetStepsByWorkflowID(ctx context.Context, workflowID uint) ([]entity.ApprovalStep, error)
	ResetStepsByWorkflowID(ctx context.Context, workflowID uint) error
}
