package gormAdapter

import (
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/filter/sql"
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

// // Filter aplica el AST de filtros a una consulta GORM
// func Filter(filter filter.Node) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if filter == nil {
// 			return db
// 		}

// 		// Inicializamos nuestro Visitor inyectándole la instancia actual de DB
// 		visitor := &GormVisitor{DB: db}

// 		if err := filter.Accept(visitor); err != nil {
// 			// Si el AST es inválido a nivel de SQL, abortamos la consulta
// 			_ = db.AddError(err)
// 			return db
// 		}

// 		// Retornamos la base de datos con todas las cláusulas Where/Or aplicadas
// 		return visitor.DB
// 	}
// }

// Filter aplica el AST de filtros a una consulta GORM
func Filter(filter filter.Clause) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return applyFilter(db, filter)
	}
}

// GFilter aplica el AST de filtros a una consulta GORM en estilo Builder
func GFilter(filter filter.Clause) func(stmt *gorm.Statement) {
	return func(stmt *gorm.Statement) {
		stmt.DB = applyFilter(stmt.DB, filter)
	}
}

func applyFilter(db *gorm.DB, clause filter.Clause) *gorm.DB {
	if clause == nil {
		return db
	}

	sql_builder := sql.NewSQLBuilder(sql.WithContext(db.Statement.Context))

	sql, args, err := sql_builder.Build(clause)
	if err != nil {
		_ = db.AddError(err)
		return db
	}

	return db.Where(sql, args...)
}
