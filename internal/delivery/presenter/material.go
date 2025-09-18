package presenter

type MaterialRequest struct {
	MaterialGroup string `json:"material_group" validate:"required,not_blank"`
	MaterialID    uint   `json:"material_id" validate:"required"`
	Active        bool   `json:"active"`
}

type MaterialResponse struct {
	ID            int    `json:"id"`
	MaterialGroup string `json:"material_group"`
	SupplierID    uint   `json:"supplier_id"`
	MaterialID    uint   `json:"material_id"`
	Active        bool   `json:"active"`
}
