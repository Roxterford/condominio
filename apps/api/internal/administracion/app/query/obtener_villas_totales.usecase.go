package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/villa"
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerVillasTotales usecase.Handler[cc.BaseContext, struct{}, *villa.VillasTotales]

type obtenerVillasTotales struct {
	finder villa.VillasTotalesFinder
}

func NewObtenerVillasTotales(finder villa.VillasTotalesFinder) ObtenerVillasTotales {
	return &obtenerVillasTotales{
		finder: finder,
	}
}

func (u *obtenerVillasTotales) Exec(
	ctx cc.BaseContext,
	input struct{},
) (*villa.VillasTotales, core.Error) {
	result, err := u.finder.ObtenerVillasTotales(ctx)

	if err != nil {
		return nil, err
	}

	return result, nil
}
