package entity

import (
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type Group struct {
	ID        uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"column:name;not null"`
	Value     string     `json:"value" gorm:"column:value;uniqueIndex"`
	IsActive  bool       `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (r *Group) TableName() string {
	return "groups"
}

func (r *Group) ToPresenter() *presenter.GroupResponse {
	return &presenter.GroupResponse{
		ID:        int(r.ID),
		GroupName: r.Name,
		Value:     r.Value,
		Active:    r.IsActive,
	}
}
