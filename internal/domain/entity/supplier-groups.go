package entity

import "time"

type SupplierGroup struct {
	SupplierID uint      `json:"supplier_id" gorm:"column:supplier_id;primaryKey"`
	GroupID    uint      `json:"group_id" gorm:"column:group_id;primaryKey"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (SupplierGroup) TableName() string {
	return "supplier_groups"
}
