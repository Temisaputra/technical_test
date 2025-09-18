package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type Address struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SupplierID uint       `json:"supplier_id" gorm:"column:supplier_id;not null;index"`
	Name       string     `json:"name" gorm:"column:name;not null"`
	Address    string     `json:"address" gorm:"column:address;not null"`
	Main       bool       `json:"main" gorm:"column:main;default:false"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Address) TableName() string {
	return "addresses"
}

func (a *Address) ToPresenter() *presenter.AddressResponse {
	return &presenter.AddressResponse{
		ID:         int(a.ID),
		SupplierID: a.SupplierID,
		Name:       a.Name,
		Address:    a.Address,
		Main:       a.Main,
	}
}
