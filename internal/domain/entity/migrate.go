package entity

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Users{},
		&Address{},
		&Supplier{},
		&Contact{},
		&Group{},
		&SupplierGroup{},
		&MaterialList{},
		&StatusHistory{},
		&ApprovalWorkflow{},
		&ApprovalLog{},
	)
}

func Drop(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&Users{},
		&Supplier{},
		&Address{},
		&Contact{},
		&Group{},
		&SupplierGroup{},
		&MaterialList{},
		&StatusHistory{},
		&ApprovalWorkflow{},
		&ApprovalLog{},
	)
}
