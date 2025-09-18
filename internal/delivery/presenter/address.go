package presenter

type AddressRequest struct {
	Name    string `json:"name" validate:"required,not_blank"`
	Address string `json:"address" validate:"required,not_blank"`
	Main    bool   `json:"main"`
}

type AddressResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	SupplierID uint   `json:"supplier_id"`
	Main       bool   `json:"main"`
}
