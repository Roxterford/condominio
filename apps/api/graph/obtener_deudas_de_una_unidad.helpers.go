package graph

import (
	"context"

	"github.com/Sanaruca/condominio/graph/model"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/unidades/app/query"
)

func (r *queryResolver) resolverDeudasBase(
	ctx context.Context,
	searchID string,
	paginator *model.Paginator,
	handler query.ObtenerDeudasDeUnaUnidadPorID,
) (*model.PaginatedDeuda, error) {
	bc, err := cc.Wrap(ctx).AsBase()
	if err != nil {
		return nil, err
	}

	res, err := handler.Exec(bc, query.ObtenerDeudasDeUnaUnidadDTO{
		SearchID:  searchID,
		Paginator: paginator.ToDomainPaginator(),
	})
	if err != nil {
		return nil, err
	}

	data := make([]*model.Deuda, len(res.Data))
	for i, d := range res.Data {
		// TODO: change
		data[i] = model.DeudaFromDomain(d, nil)
	}

	return &model.PaginatedDeuda{
		Data:  data,
		Total: int32(res.Total),
		Page:  int32(res.Page),
		Pages: int32(res.Pages),
		Limit: int32(res.Limit),
	}, nil
}
