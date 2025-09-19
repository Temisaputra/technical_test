package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type ApprovalWorkflow struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	SupplierID  uint          `json:"supplier_id"`
	Stage       string        `json:"stage"`
	SLAHours    int           `json:"sla_hours"`
	Status      string        `json:"status"`
	CurrentStep int           `gorm:"default:0" json:"current_step"` //
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Logs        []ApprovalLog `json:"logs" gorm:"foreignKey:WorkflowID"`
}

func (ApprovalWorkflow) TableName() string {
	return "approval_workflows"
}

func (a *ApprovalWorkflow) ToPresenter() *presenter.ApprovalResponse {
	return &presenter.ApprovalResponse{
		ID:         a.ID,
		SupplierID: a.SupplierID,
		Stage:      a.Stage,
		SLAHours:   a.SLAHours,
		Status:     a.Status,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

type ApprovalLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	WorkflowID uint      `json:"workflow_id"`
	Role       string    `json:"role"`
	UserName   string    `json:"user"`
	Action     string    `json:"action"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ApprovalLog) TableName() string {
	return "approval_logs"
}

func (a *ApprovalLog) ToPresenter() *presenter.ApprovalLogResponse {
	return &presenter.ApprovalLogResponse{
		ID:         a.ID,
		WorkflowID: a.WorkflowID,
		Role:       a.Role,
		UserName:   a.UserName,
		Action:     a.Action,
		Notes:      a.Notes,
		CreatedAt:  a.CreatedAt,
	}
}

type ApprovalStep struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkflowID uint      `gorm:"not null;index" json:"workflow_id"`
	StepOrder  int       `gorm:"not null" json:"step_order"`
	Role       string    `gorm:"type:varchar(100);not null" json:"role"`
	SLAHours   int       `gorm:"default:0" json:"sla_hours"`
	Status     string    `gorm:"type:varchar(50);default:'Pending'" json:"status"` // Pending, In Progress, Completed, Rejected
	AssignedTo *uint     `json:"assigned_to,omitempty"`
	IsCurrent  bool      `gorm:"default:false" json:"is_current"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ApprovalStep) TableName() string {
	return "approval_steps"
}
