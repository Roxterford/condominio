package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core/common"
)

// gormUoW implementa core.UnitOfWork[T] usando GORM
type gormUoW[T any] struct {
	db      *gorm.DB
	factory func(tx *gorm.DB) T // Sabe cómo instanciar los repositorios con el tx de GORM
}

// NewGormUnitOfWork es el constructor del UoW genérico para GORM
func NewGormUnitOfWork[T any](db *gorm.DB, factory func(tx *gorm.DB) T) common.UnitOfWork[T] {
	return &gormUoW[T]{
		db:      db,
		factory: factory,
	}
}

func (u *gormUoW[T]) Do(ctx context.Context, fn func(deps T) error) error {
	// .WithContext(ctx) asegura el rastreo, timeouts y propagación del contexto en GORM
	// .Transaction() abre la TX, maneja panics, commits y rollbacks de forma nativa
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Creamos los repositorios usando la instancia de transacción activa 'tx'
		repos := u.factory(tx)

		// 2. Ejecutamos la lógica de negocio (el caso de uso)
		// Si 'fn' retorna un error, GORM intercepta y ejecuta Rollback automáticamente
		return fn(repos)
	})
}
