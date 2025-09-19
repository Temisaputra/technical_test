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

func (u *ApprovalUsecase) GetAllApprovalWorkflows(ctx context.Context, pagination *request.Pagination, params *presenter.ApprovalRequest) ([]*presenter.ApprovalResponse, response.Meta, error) {

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

func (u *ApprovalUsecase) GetActiveWorkflowByID(ctx context.Context, Id uint) (*presenter.ApprovalResponse, error) {

	workflow, err := u.approvalRepo.GetActiveWorkflowByID(ctx, Id)
	if err != nil {
		return nil, fmt.Errorf("[ApprovalUsecase-GetActiveWorkflowByID] error: %w", err)
	}

	return workflow.ToPresenter(), nil
}

func (u *ApprovalUsecase) CreateWorkflow(ctx context.Context, workflow *presenter.ApprovalCreateRequest) error {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		currentTime := time.Now()
		storedWorkflow := &entity.ApprovalWorkflow{
			SupplierID: workflow.SupplierID,
			Status:     workflow.Status,
			Stage:      workflow.Stage,
			CreatedAt:  currentTime,
			UpdatedAt:  currentTime,
			SLAHours:   workflow.SLAHours,
		}

		workflowId, err := u.approvalRepo.CreateWorkflow(txCtx, storedWorkflow)
		if err != nil {
			return fmt.Errorf("[ApprovalUsecase-CreateWorkflow] error: %w", err)
		}

		// Create associated steps
		var steps []entity.ApprovalStep
		for _, stepReq := range workflow.Steps {
			steps = append(steps, entity.ApprovalStep{
				WorkflowID: workflowId,
				StepOrder:  stepReq.StepOrder,
				Role:       stepReq.Role,
				SLAHours:   stepReq.SLAHours,
				AssignedTo: stepReq.AssignedTo,
				Status:     "Draft",
				CreatedAt:  currentTime,
				UpdatedAt:  currentTime,
			})
		}
		if err := u.approvalRepo.CreateWorkflowSteps(txCtx, steps); err != nil {
			return fmt.Errorf("[ApprovalUsecase-CreateWorkflow] error creating steps: %w", err)
		}

		err = u.approvalRepo.CreateLog(txCtx, &entity.ApprovalLog{
			WorkflowID: workflowId,
			Action:     "Created",
			Role:       "System",
			UserName:   "System",
			Notes:      "Approval workflow created",
			CreatedAt:  currentTime,
		})
		if err != nil {
			return fmt.Errorf("[ApprovalUsecase-CreateWorkflow] error creating log: %w", err)
		}

		return nil
	})
}

func (u *ApprovalUsecase) ApproveWorkflow(ctx context.Context, req *presenter.ApprovalUpdateRequest) error {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Ambil workflow
		workflow, err := u.approvalRepo.GetActiveWorkflowByID(txCtx, req.ID)
		if err != nil {
			return fmt.Errorf("[ApproveWorkflow] error fetching workflow: %w", err)
		}
		if workflow == nil {
			return fmt.Errorf("workflow not found")
		}

		// 2. Ambil steps
		steps, err := u.approvalRepo.GetStepsByWorkflowID(txCtx, workflow.ID)
		if err != nil {
			return fmt.Errorf("[ApproveWorkflow] error fetching steps: %w", err)
		}
		if len(steps) == 0 {
			return fmt.Errorf("workflow has no steps defined")
		}

		// 3. Cari current step
		var currentStep *entity.ApprovalStep
		for i := range steps {
			if steps[i].StepOrder == workflow.CurrentStep {
				currentStep = &steps[i]
				break
			}
		}
		if currentStep == nil {
			return fmt.Errorf("no current step found for workflow")
		}

		// 4. Update current step → Done
		currentStep.Status = "Completed"
		currentStep.IsCurrent = false
		if err := u.approvalRepo.UpdateWorkflowStep(txCtx, currentStep); err != nil {
			return fmt.Errorf("[ApproveWorkflow] error updating current step: %w", err)
		}

		// 5. Tambahkan log
		now := time.Now()
		err = u.approvalRepo.CreateLog(txCtx, &entity.ApprovalLog{
			WorkflowID: workflow.ID,
			Action:     req.Reviewer.Action,
			Role:       req.Reviewer.Role,
			UserName:   req.Reviewer.UserName,
			Notes:      req.Reviewer.Notes,
			CreatedAt:  now,
		})
		if err != nil {
			return fmt.Errorf("[ApproveWorkflow] error creating log: %w", err)
		}

		// 6. Cek apakah masih ada step berikutnya
		if workflow.CurrentStep < len(steps) {
			workflow.CurrentStep++ // pindah ke step berikutnya

			// update step berikutnya jadi In Progress
			for i := range steps {
				if steps[i].StepOrder == workflow.CurrentStep {
					steps[i].Status = "In Progress"
					steps[i].IsCurrent = true
					if err := u.approvalRepo.UpdateWorkflowStep(txCtx, &steps[i]); err != nil {
						return fmt.Errorf("[ApproveWorkflow] error updating next step: %w", err)
					}
					break
				}
			}

			// update workflow (masih ongoing)
			workflow.Status = "In Progress"

			switch workflow.Stage {
			case "Draft":
				workflow.Stage = "In Review"
			case "In Review":
				workflow.Stage = "In Assessment"
			case "In Assessment":
				workflow.Stage = "Final Review"
			}

			workflow.UpdatedAt = now
			if err := u.approvalRepo.UpdateWorkflow(txCtx, workflow); err != nil {
				return fmt.Errorf("[ApproveWorkflow] error updating workflow: %w", err)
			}
		} else {
			// kalau step terakhir selesai → workflow Completed
			workflow.Status = "Completed"
			workflow.Stage = "Active"
			if err := u.approvalRepo.UpdateWorkflow(txCtx, workflow); err != nil {
				return fmt.Errorf("[ApproveWorkflow] error completing workflow: %w", err)
			}
		}

		return nil
	})
}

func (u *ApprovalUsecase) GetLogsByWorkflowID(ctx context.Context, workflowID uint) ([]presenter.ApprovalLogResponse, error) {
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

func (u *ApprovalUsecase) CreateLog(ctx context.Context, log *presenter.ApprovalLogCreateRequest) error {
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
