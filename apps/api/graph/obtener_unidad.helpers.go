package graph

import (
	"context"

	"github.com/Sanaruca/condominio/graph/model"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/unidades/app/query"
)

func (r *queryResolver) resolverUnidadBase(
	ctx context.Context,
	searchID string,
	handler query.ObtenerUnidad, // O el tipo de interfaz que compartan
) (*model.Unidad, error) {

	bc, err := cc.Wrap(ctx).AsBase()
	if err != nil {
		return nil, err
	}

	res, err := handler.Exec(bc, query.ObtenerUnidadDTO{SearchID: searchID})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return model.UnidadFromDomain(*res), nil
}
