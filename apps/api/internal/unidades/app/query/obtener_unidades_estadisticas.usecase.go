package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type ObtenerUnidadesEstadisticas usecase.Handler[cc.BaseContext, any, *unidad.Estadisticas]

type obtenerUnidadesEstadisticas struct {
	finder unidad.EstadisticasFinder
}

func NewObtenerUnidadesEstadisticas(finder unidad.EstadisticasFinder) ObtenerUnidadesEstadisticas {
	return &obtenerUnidadesEstadisticas{
		finder: finder,
	}
}

func (u *obtenerUnidadesEstadisticas) Exec(
	ctx cc.BaseContext,
	input any,
) (*unidad.Estadisticas, core.Error) {

	result, err := u.finder.Obtener(ctx)

	if err != nil {
		return nil, err
	}

	return result, nil
}
