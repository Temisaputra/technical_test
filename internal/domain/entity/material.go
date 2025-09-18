package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type MaterialList struct {
	ID            uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SupplierID    uint       `json:"supplier_id" gorm:"column:supplier_id;not null"`
	MaterialGroup string     `json:"material_group" gorm:"column:material_group"`
	MaterialID    uint       `json:"material_id" gorm:"column:material_id"`
	Description   string     `json:"description" gorm:"column:description"`
	Price         float64    `json:"price" gorm:"column:price;type:numeric(18,2);default:0"`
	IsActive      bool       `json:"is_active" gorm:"column:is_active"`
	CreatedAt     time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (MaterialList) TableName() string {
	return "material_lists"
}

func (m *MaterialList) ToPresenter() *presenter.MaterialResponse {
	return &presenter.MaterialResponse{
		ID:            int(m.ID),
		SupplierID:    m.SupplierID,
		MaterialGroup: m.MaterialGroup,
		MaterialID:    m.MaterialID,
		Active:        m.IsActive,
	}
}
