package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/pkg/helper"
	"github.com/gorilla/mux"
)

type approvalUsecase interface {
	GetAllApprovalWorkflows(ctx context.Context, pagination *request.Pagination, params *presenter.ApprovalRequest) ([]*presenter.ApprovalResponse, response.Meta, error)
	GetActiveWorkflowByID(ctx context.Context, Id uint) (*presenter.ApprovalResponse, error)
	CreateWorkflow(ctx context.Context, workflow *presenter.ApprovalCreateRequest) error
	ApproveWorkflow(ctx context.Context, workflow *presenter.ApprovalUpdateRequest) error
	GetLogsByWorkflowID(ctx context.Context, workflowID uint) ([]presenter.ApprovalLogResponse, error)
	CreateLog(ctx context.Context, log *presenter.ApprovalLogCreateRequest) error
}

type ApprovalHandler struct {
	approvalUsecase approvalUsecase
}

func NewApprovalHandler(uc approvalUsecase) *ApprovalHandler {
	return &ApprovalHandler{approvalUsecase: uc}
}

// @Tags Approval
// @Summary Get All Approval Workflows
// @Description Get All Approval Workflows
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param stage query string false "Filter by stage"
// @Param supplier_id query int false "Filter by supplier ID"
// @Success 200 {object} helper.Response{data=[]presenter.ApprovalResponse,meta=response.Meta}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /approvals [get]
func (h *ApprovalHandler) GetAllApprovalWorkflows(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	pagination := &request.Pagination{
		OrderBy:   r.URL.Query().Get("order_by"),
		OrderType: r.URL.Query().Get("order_type"),
		Page:      page,
		PageSize:  pageSize,
	}

	params := &presenter.ApprovalRequest{
		Status: r.URL.Query().Get("status"),
		Stage:  r.URL.Query().Get("stage"),
	}

	if sid := r.URL.Query().Get("supplier_id"); sid != "" {
		if val, err := strconv.Atoi(sid); err == nil {
			params.SupplierID = uint(val)
		}
	}

	data, meta, err := h.approvalUsecase.GetAllApprovalWorkflows(r.Context(), pagination, params)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
		Meta:       &meta,
	})
}

// @Tags Approval
// @Summary Get Active Workflow By Supplier ID
// @Description Get Active Workflow By Supplier ID
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} helper.Response{data=presenter.ApprovalResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /approval/{id} [get]
func (h *ApprovalHandler) GetActiveWorkflowByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}
	data, err := h.approvalUsecase.GetActiveWorkflowByID(r.Context(), uint(idInt))
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
	})
}

// @Tags Approval
// @Summary Create Workflow
// @Description Create a new workflow for supplier
// @Accept json
// @Produce json
// @Param request body presenter.ApprovalCreateRequest true "Workflow data"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /create-workflow [post]
func (h *ApprovalHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var params presenter.ApprovalCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.approvalUsecase.CreateWorkflow(r.Context(), &params); err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusCreated,
		Message:    "success",
	})
}

// @Tags Approval
// @Summary Create Approval Log
// @Description Submit approval action (approve, reject, comment)
// @Accept json
// @Produce json
// @Param id path int true "Workflow ID"
// @Param request body presenter.ApprovalLogCreateRequest true "Approval log data"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /approvals/logs/{id} [post]
func (h *ApprovalHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("workflow id is required"), nil)
		return
	}

	var params presenter.ApprovalLogCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params.WorkflowID = uint(idInt)

	if err := h.approvalUsecase.CreateLog(r.Context(), &params); err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusCreated,
		Message:    "success",
	})
}

// @Tags Approval
// @Summary Get Logs By Workflow ID
// @Description Get all logs of a workflow
// @Accept json
// @Produce json
// @Param id path int true "Workflow ID"
// @Success 200 {object} helper.Response{data=[]presenter.ApprovalLogResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /approvals/{id}/logs [get]
func (h *ApprovalHandler) GetLogsByWorkflowID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("workflow id is required"), nil)
		return
	}

	data, err := h.approvalUsecase.GetLogsByWorkflowID(r.Context(), uint(idInt))
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
	})
}

// @Tags Approval
// @Summary Approve Workflow
// @Description Approve an existing workflow
// @Accept json
// @Produce json
// @Param id path int true "Workflow ID"
// @Param request body presenter.ApprovalUpdateRequest true "Approve workflow data"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /approve-workflow/{id} [put]
func (h *ApprovalHandler) ApproveWorkflow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("workflow id is required"), nil)
		return
	}
	var params presenter.ApprovalUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params.ID = uint(idInt)
	if err := h.approvalUsecase.ApproveWorkflow(r.Context(), &params); err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	helper.WriteResponse(w, nil, &helper.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
	})
}
