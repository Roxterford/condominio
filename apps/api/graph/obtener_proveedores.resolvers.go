package graph

import (
	"context"

	"github.com/Sanaruca/condominio/graph/model"
	corecontext "github.com/Sanaruca/condominio/internal/core/context"
)

func (r *queryResolver) ObtenerProveedores(ctx context.Context) ([]*model.Proveedor, error) {
	// TODO: Cambiar a AsAdmin cuando la autenticación esté funcional
	adminCtx, err := corecontext.Wrap(ctx).AsBase()
	if err != nil {
		return nil, err
	}

	proveedores, err := r.Administracion.Queries.ObtenerProveedores.Exec(adminCtx, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Proveedor, len(proveedores))
	for i, p := range proveedores {
		email := p.Email().String()
		telefono := p.Telefono().String()
		result[i] = &model.Proveedor{
			ID:            p.ID(),
			Rif:           p.Rif().String(),
			Nombre:        p.Nombre(),
			Tipo:          string(p.Tipo()),
			Email:         &email,
			Telefono:      &telefono,
			Direccion:     p.Direccion(),
			CreadoEn:      p.CreadoEn(),
			ActualizadoEn: p.ActualizadoEn(),
		}
	}

	return result, nil
}
