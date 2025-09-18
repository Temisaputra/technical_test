package presenter

import "time"

// ApprovalResponse adalah DTO untuk mengembalikan data ke client
type ApprovalResponse struct {
	ID         uint      `json:"id"`
	SupplierID uint      `json:"supplier_id"`
	Stage      string    `json:"stage"`
	SLAHours   int       `json:"sla_hours"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ApprovalRequest struct {
	Status     string `json:"status"`
	Stage      string `json:"stage"`
	SupplierID uint   `json:"supplier_id"`
}

// Untuk membuat workflow baru
type ApprovalCreateRequest struct {
	SupplierID uint   `json:"supplier_id" validate:"required"`
	Stage      string `json:"stage" validate:"required"` // Draft, In Review, In Assessment, Active
	SLAHours   int    `json:"sla_hours"`
	Status     string `json:"status"` // default: In Progress
}

// Untuk update workflow (misal ganti stage/status)
type ApprovalUpdateRequest struct {
	ID         uint   `json:"id" validate:"required"`
	Stage      string `json:"stage"`
	Status     string `json:"status"`
	SLAHours   int    `json:"sla_hours"`
	SupplierID uint   `json:"supplier_id" validate:"required"`
}

// Untuk membuat log baru
type ApprovalLogCreateRequest struct {
	WorkflowID uint   `json:"workflow_id"`
	Role       string `json:"role" validate:"required"`
	UserName   string `json:"user" validate:"required"`
	Action     string `json:"action" validate:"required"` // Approve, Return, Comment
	Notes      string `json:"notes"`
}

type ApprovalLogResponse struct {
	ID         uint      `json:"id"`
	WorkflowID uint      `json:"workflow_id"`
	Role       string    `json:"role"`
	UserName   string    `json:"user"`
	Action     string    `json:"action"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}
