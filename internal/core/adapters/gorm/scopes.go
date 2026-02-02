package gormAdapter

import (
	"github.com/Sanaruca/condominio/internal/core/common"
	"gorm.io/gorm"
)

func Paginate(paginator common.Paginator) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		paginator.Sanitize()
		offset := (paginator.Page - 1) * paginator.Limit
		return db.Offset(offset).Limit(paginator.Limit)
	}
}
func GPaginate(paginator common.Paginator) func(*gorm.Statement) {
	return func(stmt *gorm.Statement) {
		paginator.Sanitize()
		offset := (paginator.Page - 1) * paginator.Limit
		stmt.DB.Offset(offset).Limit(paginator.Limit)
	}
}
