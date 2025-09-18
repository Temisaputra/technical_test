package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/validation"
	"github.com/Temisaputra/warOnk/pkg/helper"
	"github.com/gorilla/mux"
)

type supplierUsecase interface {
	GetAllSupplier(ctx context.Context, pagination *request.Pagination, params *presenter.SupplierRequest) (suppliers []*presenter.SupplierResponse, meta response.Meta, err error)
	GetSupplierByID(ctx context.Context, id int) (supplier *presenter.SupplierResponse, err error)
	CreateSupplier(ctx context.Context, params *presenter.SupplierCreateRequest) error
	UpdateSupplier(ctx context.Context, params *presenter.SupplierCreateRequest, id int) error
	DeleteSupplier(ctx context.Context, id int) error
	BlockUnblockSupplier(ctx context.Context, id int, params presenter.SupplierBlockUnblockRequest) error
	ExportSupplier(ctx context.Context, params *presenter.SupplierRequest) (fileName string, err error)
}

type SupplierHandler struct {
	supplierUsecase supplierUsecase
}

func NewSupplierHandler(supplierUsecase supplierUsecase) *SupplierHandler {
	return &SupplierHandler{
		supplierUsecase: supplierUsecase,
	}
}

// GetAllSupplier godoc
// @Tags Supplier
// @Summary Get All Supplier
// @Description Get All Supplier
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
// @Success 200 {object} helper.Response{data=[]presenter.SupplierResponse,meta=response.Meta}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /suppliers [get]
func (h *SupplierHandler) GetAllSupplier(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	pagination := &request.Pagination{
		Keyword:   r.URL.Query().Get("keyword"),
		OrderBy:   r.URL.Query().Get("order_by"),
		OrderType: r.URL.Query().Get("order_type"),
		Page:      page,
		PageSize:  pageSize,
	}

	params := &presenter.SupplierRequest{
		Name:     r.URL.Query().Get("name"),
		NickName: r.URL.Query().Get("nick_name"),
		Status:   r.URL.Query().Get("status"),
	}

	data, meta, err := h.supplierUsecase.GetAllSupplier(r.Context(), pagination, params)
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

// GetSupplierByID godoc
// @Tags Supplier
// @Summary Get Supplier by ID
// @Description Get Supplier by ID
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} helper.Response{data=presenter.SupplierResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier/{id} [get]
func (h *SupplierHandler) GetSupplierByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}
	data, err := h.supplierUsecase.GetSupplierByID(r.Context(), idInt)
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

// CreateSupplier godoc
// @Tags Supplier
// @Summary Create a new supplier
// @Description Create a new supplier
// @Accept json
// @Produce json
// @Param request body presenter.SupplierCreateRequest true "Supplier data"
// @Success 201 {object} helper.Response{data=presenter.SupplierResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier-create [post]
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var params presenter.SupplierCreateRequest
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

	err = h.supplierUsecase.CreateSupplier(r.Context(), &params)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusCreated
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}

// UpdateSupplier godoc
// @Tags Supplier
// @Summary Update a supplier
// @Description Update a supplier
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Param request body presenter.SupplierCreateRequest true "Supplier data"
// @Success 200 {object} helper.Response{data=presenter.SupplierResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier-update/{id} [put]
func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}

	var params presenter.SupplierCreateRequest
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

	err = h.supplierUsecase.UpdateSupplier(r.Context(), &params, idInt)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}

// DeleteSupplier godoc
// @Tags Supplier
// @Summary Delete a supplier
// @Description Delete a supplier
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} helper.Response{data=presenter.SupplierResponse}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier-delete/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}

	err := h.supplierUsecase.DeleteSupplier(r.Context(), idInt)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}

	var response helper.Response

	response.StatusCode = http.StatusOK
	response.Message = "success"

	helper.WriteResponse(w, nil, &response)
}

// BlockUnblockSupplier godoc
// @Tags Supplier
// @Summary Block Unblock a supplier
// @Description Block Unblock a supplier
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Param request body presenter.SupplierBlockUnblockRequest true "Supplier data"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier-block-unblock/{id} [patch]
func (h *SupplierHandler) BlockUnblockSupplier(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)
	if idInt == 0 {
		helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
		return
	}
	var params presenter.SupplierBlockUnblockRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.supplierUsecase.BlockUnblockSupplier(r.Context(), idInt, params)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}
	var response helper.Response
	response.StatusCode = http.StatusOK
	response.Message = "success"
	helper.WriteResponse(w, nil, &response)
}

// ExportSupplier godoc
// @Tags Supplier
// @Summary Export Supplier data to Excel
// @Description Export Supplier data to Excel
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword for search"
// @Param name query string false "Filter by name"
// @Param nick_name query string false "Filter by nick_name"
// @Param status query string false "Filter by status"
// @Param address query string false "Filter by address"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /supplier-export [get]
func (h *SupplierHandler) ExportSupplier(w http.ResponseWriter, r *http.Request) {
	params := &presenter.SupplierRequest{
		Name:     r.URL.Query().Get("name"),
		NickName: r.URL.Query().Get("nick_name"),
		Status:   r.URL.Query().Get("status"),
	}
	fileName, err := h.supplierUsecase.ExportSupplier(r.Context(), params)
	if err != nil {
		helper.WriteResponse(w, err, nil)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	file.Close()

	err = os.Remove(fileName)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}
}
