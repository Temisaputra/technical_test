package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/delivery/repository"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
)

type ApprovalUsecase struct {
	approvalRepo    repository.ApprovalRepository
	transactionRepo repository.TransactionRepository
}

func NewApprovalUsecase(
	approvalRepository repository.ApprovalRepository,
	transactionRepository repository.TransactionRepository,
) *ApprovalUsecase {
	return &ApprovalUsecase{
		approvalRepo:    approvalRepository,
		transactionRepo: transactionRepository,
	}
}

func (u *ApprovalUsecase) GetAllApprovalWorkflows(
	ctx context.Context,
	pagination *request.Pagination,
	params *presenter.ApprovalRequest,
) ([]*presenter.ApprovalResponse, response.Meta, error) {

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	res, meta, err := u.approvalRepo.GetAllApprovalWorkflows(ctx, pagination, params)
	if err != nil {
		return nil, meta, fmt.Errorf("[ApprovalUsecase-GetAllApprovalWorkflows] error: %w", err)
	}

	return res, meta, nil
}

func (u *ApprovalUsecase) GetActiveWorkflowBySupplierID(
	ctx context.Context,
	supplierID uint,
) (*presenter.ApprovalResponse, error) {

	workflow, err := u.approvalRepo.GetActiveWorkflowBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("[ApprovalUsecase-GetActiveWorkflowBySupplierID] error: %w", err)
	}

	return workflow.ToPresenter(), nil
}

func (u *ApprovalUsecase) CreateWorkflow(
	ctx context.Context,
	workflow *presenter.ApprovalCreateRequest,
) error {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		currentTime := time.Now()
		workflow := &entity.ApprovalWorkflow{
			SupplierID: workflow.SupplierID,
			Status:     workflow.Status,
			Stage:      workflow.Stage,
			CreatedAt:  currentTime,
			UpdatedAt:  currentTime,
			SLAHours:   workflow.SLAHours,
		}

		if err := u.approvalRepo.CreateWorkflow(txCtx, workflow); err != nil {
			return fmt.Errorf("[ApprovalUsecase-CreateWorkflow] error: %w", err)
		}
		return nil
	})
}

func (u *ApprovalUsecase) UpdateWorkflow(
	ctx context.Context,
	workflow *presenter.ApprovalUpdateRequest,
) error {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		currentTime := time.Now()
		workflow := &entity.ApprovalWorkflow{
			ID:         workflow.ID,
			SupplierID: workflow.SupplierID,
			Status:     workflow.Status,
			Stage:      workflow.Stage,
			UpdatedAt:  currentTime,
			SLAHours:   workflow.SLAHours,
		}
		if err := u.approvalRepo.UpdateWorkflow(txCtx, workflow); err != nil {
			return fmt.Errorf("[ApprovalUsecase-UpdateWorkflow] error: %w", err)
		}
		return nil
	})
}

func (u *ApprovalUsecase) GetLogsByWorkflowID(
	ctx context.Context,
	workflowID uint,
) ([]presenter.ApprovalLogResponse, error) {
	res := []presenter.ApprovalLogResponse{}
	logs, err := u.approvalRepo.GetLogsByWorkflowID(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("[ApprovalUsecase-GetLogsByWorkflowID] error: %w", err)
	}

	for _, log := range logs {
		res = append(res, *log.ToPresenter())
	}
	return res, nil
}

func (u *ApprovalUsecase) CreateLog(
	ctx context.Context,
	log *presenter.ApprovalLogCreateRequest,
) error {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		currentTime := time.Now()
		log := &entity.ApprovalLog{
			WorkflowID: log.WorkflowID,
			Action:     log.Action,
			Role:       log.Role,
			UserName:   log.UserName,
			Notes:      log.Notes,
			CreatedAt:  currentTime,
		}

		if err := u.approvalRepo.CreateLog(txCtx, log); err != nil {
			return fmt.Errorf("[ApprovalUsecase-CreateLog] error: %w", err)
		}
		return nil
	})
}
