package gorm

import (
	"context"
	"errors"

	proveedor_pkg "github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"gorm.io/gorm"
)

type GORMProveedorRepository struct {
	db      *gorm.DB
	factory *proveedor_pkg.ProveedorFactory
}

func NewGORMProveedorRepository(
	db *gorm.DB,
	factory *proveedor_pkg.ProveedorFactory,
) proveedor_pkg.ProveedorRepository {
	return &GORMProveedorRepository{db: db, factory: factory}
}

func (r *GORMProveedorRepository) Guardar(
	ctx context.Context,
	proveedor *proveedor_pkg.Proveedor,
) core.Error {
	model := toProveedorTable(proveedor)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return proveedor_pkg.ErrProveedorDuplicado
		}
		return core.WrapError(err)
	}
	return nil
}

func (r *GORMProveedorRepository) ObtenerPorID(
	ctx context.Context,
	id string,
) (*proveedor_pkg.Proveedor, core.Error) {
	var model Proveedor
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proveedor_pkg.ErrProveedorNoEncontrado
		}
		return nil, core.WrapError(err)
	}

	return r.toProveedor(&model), nil
}

func (r *GORMProveedorRepository) ObtenerTodos(
	ctx context.Context,
) ([]*proveedor_pkg.Proveedor, core.Error) {
	var models []Proveedor
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, core.WrapError(err)
	}
	proveedores := make([]*proveedor_pkg.Proveedor, len(models))
	for i, m := range models {
		proveedores[i] = r.factory.Assemble(
			m.ID,
			m.Rif,
			m.Nombre,
			m.Email,
			m.Telefono,
			m.Direccion,
			m.Registro,
			m.Actualizacion,
		)
	}
	return proveedores, nil
}

func (r *GORMProveedorRepository) Actualizar(
	ctx context.Context,
	proveedor *proveedor_pkg.Proveedor,
) core.Error {
	model := toProveedorTable(proveedor)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return core.WrapError(err)
	}
	return nil
}

func (r *GORMProveedorRepository) Eliminar(ctx context.Context, id string) core.Error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Proveedor{}).Error; err != nil {
		return core.WrapError(err)
	}
	return nil
}

func (r *GORMProveedorRepository) ExistePorRif(ctx context.Context, rif string) (bool, core.Error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Proveedor{}).Where("rif = ?", rif).Count(&count).Error; err != nil {
		return false, core.WrapError(err)
	}
	return count > 0, nil
}

func (r *GORMProveedorRepository) ExistePorID(ctx context.Context, id string) (bool, core.Error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Proveedor{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, core.WrapError(err)
	}
	return count > 0, nil
}

func toProveedorTable(p *proveedor_pkg.Proveedor) *Proveedor {
	return &Proveedor{
		ID:            p.ID(),
		Rif:           p.Rif().String(),
		Nombre:        p.Nombre(),
		Email:         p.Email().String(),
		Telefono:      p.Telefono().String(),
		Direccion:     p.Direccion(),
		Registro:      p.CreadoEn(),
		Actualizacion: p.ActualizadoEn(),
	}
}

func (r *GORMProveedorRepository) toProveedor(table *Proveedor) *proveedor_pkg.Proveedor {
	return r.factory.Assemble(
		table.ID,
		table.Rif,
		table.Nombre,
		table.Email,
		table.Telefono,
		table.Direccion,
		table.Registro,
		table.Actualizacion,
	)

}
