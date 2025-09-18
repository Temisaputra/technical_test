package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/validation"
	"github.com/Temisaputra/warOnk/pkg/helper"
	"github.com/gorilla/mux"
)

type groupUsecase interface {
	GetAllGroup(ctx context.Context, pagination *request.Pagination) (groups []*presenter.GroupResponse, meta response.Meta, err error)
	GetGroupByID(ctx context.Context, id int) (group *presenter.GroupResponse, err error)
	CreateGroup(ctx context.Context, params *presenter.GroupRequest) error
	UpdateGroup(ctx context.Context, params *presenter.GroupRequest, id int) error
	DeleteGroup(ctx context.Context, id int) error
}

type GroupHandler struct {
	groupUsecase groupUsecase
}

func NewGroupHandler(groupUsecase groupUsecase) *GroupHandler {
	return &GroupHandler{
		groupUsecase: groupUsecase,
	}
}

// GetAllGroup godoc
// @Tags Group
// @Summary Get All Group
// @Description Get All Group
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param keyword query string false "Keyword for search"
// @Param order_by query string false "Order by field"
// @Param order_type query string false "Order type (asc/desc)"
// @Param name query string false "Filter by name"
// @Param nick_name query string false "Filter by nick_name"
// @Param status query string false "Filter by status"
// @Param address query string false "Filter by address"
// @Success 200 {object} helper.Response{data=[]presenter.GroupResponse,meta=response.Meta}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /groups [get]
func (h *GroupHandler) GetAllGroup(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	pagination := &request.Pagination{
		Keyword:   r.URL.Query().Get("keyword"),
		OrderBy:   r.URL.Query().Get("order_by"),
		OrderType: r.URL.Query().Get("order_type"),
		Page:      page,
		PageSize:  pageSize,
	}

	data, meta, err := h.groupUsecase.GetAllGroup(r.Context(), pagination)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"
	response.Meta = &meta
	response.Data = data

	helper.WriteResponse(w, nil, &response)
}

// GetGroupByID godoc
// @Tags Group
// @Summary Get Group by ID
// @Description Get Group by ID
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} helper.Response{data=presenter.GroupResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /group/{id} [get]
func (h *GroupHandler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}
	data, err := h.groupUsecase.GetGroupByID(r.Context(), idInt)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"
	response.Data = data

	helper.WriteResponse(w, nil, &response)
}

// CreateGroup godoc
// @Tags Group
// @Summary Create a new group
// @Description Create a new group
// @Accept json
// @Produce json
// @Param request body presenter.GroupRequest true "Group data"
// @Success 201 {object} helper.Response{data=presenter.GroupResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /group-create [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var params presenter.GroupRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Validasi sebelum lanjut ke usecase
	if err := validation.ValidateStruct(params); err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	err = h.groupUsecase.CreateGroup(r.Context(), &params)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusCreated
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}

// UpdateGroup godoc
// @Tags Group
// @Summary Update a group
// @Description Update a group
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body presenter.GroupRequest true "Group data"
// @Success 200 {object} helper.Response{data=presenter.GroupResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /group-update/{id} [put]
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}

	var params presenter.GroupRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Validasi sebelum lanjut ke usecase
	if err := validation.ValidateStruct(params); err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	err = h.groupUsecase.UpdateGroup(r.Context(), &params, idInt)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}

// DeleteGroup godoc
// @Tags Group
// @Summary Delete a group
// @Description Delete a group
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} helper.Response{data=presenter.GroupResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /group-delete/{id} [delete]
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}

	err := h.groupUsecase.DeleteGroup(r.Context(), idInt)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}
