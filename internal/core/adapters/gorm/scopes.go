package gormAdapter

import (
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"gorm.io/gorm"
)

func Paginate(p common.Paginator) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return applyPagination(db, p) }
}

func GPaginate(p common.Paginator) func(*gorm.Statement) {
	return func(stmt *gorm.Statement) { applyPagination(stmt.DB, p) }
}

func applyPagination(db *gorm.DB, p common.Paginator) *gorm.DB {
	p.Sanitize()
	offset := p.Offset()
	return db.Offset(offset).Limit(p.Limit)
}

////////////////////////////////////////////////////////////////////////

func Filter[T filter.Filterable](filter filter.Filter[T]) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return applyFilter(db, filter) }
}

func GFilter[T filter.Filterable](filter filter.Filter[T]) func(*gorm.Statement) {
	return func(stmt *gorm.Statement) { stmt.DB = applyFilter(stmt.DB, filter) }
}

func applyFilter[T filter.Filterable](db *gorm.DB, filter filter.Filter[T]) *gorm.DB {
	// Lógica de AND
	for _, sub_filter := range filter.And {

		db = db.Session(&gorm.Session{}).Where(sub_filter)
	}

	return db
}
