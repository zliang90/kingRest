package db

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zliang90/kingRest/pkg/util/uuid"
)

type BaseModel struct {
	Id        string `gorm:"type:char(36);primary_key" json:"id"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}

func (*BaseModel) BeforeCreate(scope *gorm.Scope) error {
	if !scope.HasError() {
		// id
		IdField, hasIdField := scope.FieldByName("Id")
		if hasIdField {
			if IdField.IsBlank {
				scope.SetColumn("Id", uuid.New())
			}
		}
		// created time
		if scope.HasColumn("CreatedAt") {
			scope.SetColumn("CreateAt", time.Now())
		}
	}
	return nil
}

func (*BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	if !scope.HasError() {
		if scope.HasColumn("UpdateAt") {
			scope.SetColumn("UpdateAt", time.Now())
		}
	}
	return nil
}
