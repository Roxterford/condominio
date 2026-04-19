package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/unidades"
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerUnidadesEstadisticas usecase.Handler[cc.BaseContext, struct{}, *unidades.UnidadesEstadisticas]

type obtenerUnidadesEstadisticas struct {
	finder unidades.UnidadesEstadisticasFinder
}

func NewObtenerUnidadesEstadisticas(finder unidades.UnidadesEstadisticasFinder) ObtenerUnidadesEstadisticas {
	return &obtenerUnidadesEstadisticas{
		finder: finder,
	}
}

func (u *obtenerUnidadesEstadisticas) Exec(
	ctx cc.BaseContext,
	input struct{},
) (*unidades.UnidadesEstadisticas, core.Error) {
	result, err := u.finder.Obtener(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
