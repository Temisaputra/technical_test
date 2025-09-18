package entity

import "time"

type StatusHistory struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SupplierID uint       `json:"supplier_id" gorm:"column:supplier_id;not null"`
	Status     string     `json:"status" gorm:"column:status;not null"` // Active, Blocked, etc
	Reason     string     `json:"reason" gorm:"column:reason"`
	UpdatedBy  string     `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (StatusHistory) TableName() string {
	return "status_histories"
}
