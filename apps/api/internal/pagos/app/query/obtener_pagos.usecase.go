// Deprecated: Query legacy. Usar internal/transacciones/app/query/ en su lugar.
package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
)

type ObtenerPagosDTO struct {
	common.Paginator
	Filter *filter.Filter[pago.Pago]
}

type ObtenerPagos usecase.Handler[cc.BaseContext, ObtenerPagosDTO, *common.Paginated[pago.Pago]]

type obtenerPagos struct {
	repo pago.PagoRepository
}

func NewObtenerPagos(repo pago.PagoRepository) ObtenerPagos {

	if repo == nil {
		panic("repo is nil")
	}

	return &obtenerPagos{repo}
}

// Exec implements [ObtenerPagos].
func (o *obtenerPagos) Exec(
	ctx cc.BaseContext,
	input ObtenerPagosDTO,
) (*common.Paginated[pago.Pago], core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	ftr, e := input.Filter.Build()

	if e != nil {
		return nil, core.WrapError(e)
	}

	return o.repo.Obtener(ctx, ftr, input.Paginator)
}

func (input ObtenerPagosDTO) Validate() core.Error {
	input.Paginator.Sanitize()

	return nil

}
