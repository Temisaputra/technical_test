package presenter

import "time"

type SupplierRequest struct {
	Name     string         `json:"name" validate:"required,not_blank"`
	NickName string         `json:"nick_name" validate:"required,not_blank"`
	Status   string         `json:"status" validate:"required"`
	Address  AddressRequest `json:"address" validate:"required,not_blank"`
	Contact  ContactRequest `json:"contact" validate:"required,not_blank"`
}

type SupplierResponse struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	NickName  string             `json:"nick_name"`
	Address   AddressResponse    `json:"address"`
	Contact   ContactResponse    `json:"contact"`
	Group     GroupResponse      `json:"group"`
	Material  []MaterialResponse `json:"material"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at"`
}

type SupplierUpdateRequest struct {
	Name     string           `json:"name" validate:"required,not_blank"`
	NickName string           `json:"nick_name" validate:"required,not_blank"`
	Address  []AddressRequest `json:"address" validate:"required,not_blank"`
	Contact  []ContactRequest `json:"contact" validate:"required,not_blank"`
	Status   string           `json:"status" validate:"required"`
}

type SupplierCreateRequest struct {
	Name      string            `json:"name" validate:"required,not_blank"`
	NickName  string            `json:"nick_name" validate:"required,not_blank"`
	Address   []AddressRequest  `json:"address" validate:"required,not_blank"`
	Contact   []ContactRequest  `json:"contact" validate:"required,not_blank"`
	Groups    []uint            `json:"groups"`
	Materials []MaterialRequest `json:"materials"`
	Status    string            `json:"status" validate:"required"`
}

type SupplierBlockUnblockRequest struct {
	Reason   string `json:"reason" validate:"required,not_blank"`
	FileName string `json:"file_name"`
}
