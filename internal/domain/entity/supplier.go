package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type Supplier struct {
	ID            uint            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name          string          `json:"name" gorm:"column:name"`
	NickName      string          `json:"nick_name" gorm:"column:nick_name"`
	Status        string          `json:"status" gorm:"column:status"`
	Addresses     []Address       `json:"addresses" gorm:"foreignKey:SupplierID"`
	Contacts      []Contact       `json:"contacts" gorm:"foreignKey:SupplierID"`
	Groups        []Group         `json:"groups" gorm:"many2many:supplier_groups"`
	MaterialLists []MaterialList  `json:"material_lists" gorm:"foreignKey:SupplierID"`
	StatusHistory []StatusHistory `json:"status_history" gorm:"foreignKey:SupplierID"`
	CreatedAt     time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     *time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (r *Supplier) TableName() string {
	return "suppliers"
}

func (s *Supplier) ToPresenter() *presenter.SupplierResponse {
	return &presenter.SupplierResponse{
		ID:        int(s.ID),
		Name:      s.Name,
		NickName:  s.NickName,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
