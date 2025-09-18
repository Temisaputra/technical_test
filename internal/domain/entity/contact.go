package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type Contact struct {
	ID          uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SupplierID  uint       `json:"supplier_id" gorm:"column:supplier_id;not null;index"`
	Name        string     `json:"name" gorm:"column:name;not null"`
	JobPosition string     `json:"job_position" gorm:"column:job_position"`
	Email       string     `json:"email" gorm:"column:email"`
	Phone       string     `json:"phone" gorm:"column:phone"`
	Mobile      string     `json:"mobile" gorm:"column:mobile"`
	Main        bool       `json:"main" gorm:"column:main;default:false"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (r *Contact) TableName() string {
	return "contacts"
}

func (c *Contact) ToPresenter() *presenter.ContactResponse {
	return &presenter.ContactResponse{
		ID:          int(c.ID),
		SupplierID:  c.SupplierID,
		Name:        c.Name,
		JobPosition: c.JobPosition,
		Email:       c.Email,
		Phone:       c.Phone,
		Mobile:      c.Mobile,
		Main:        c.Main,
	}
}
