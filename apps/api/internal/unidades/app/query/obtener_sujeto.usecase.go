package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

type ObtenerSujetoDTO struct {
	SujetoID string
}

type ObtenerSujeto usecase.Handler[cc.BaseContext, ObtenerSujetoDTO, sujeto.Sujeto]

type obtenerSujeto struct {
	repo sujeto.SujetoRepository
}

func NewObtenerSujeto(repo sujeto.SujetoRepository) ObtenerSujeto {
	if repo == nil {
		panic("repository is nil")
	}
	return &obtenerSujeto{repo}
}

// Exec implements [ObtenerSujeto].
func (uc *obtenerSujeto) Exec(
	ctx cc.BaseContext,
	input ObtenerSujetoDTO,
) (sujeto.Sujeto, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	_sujeto, err := uc.repo.ObtenerPorID(ctx, sujeto.SujetoID(input.SujetoID))

	if err != nil {
		return nil, err
	}

	return _sujeto, nil
}

func (input ObtenerSujetoDTO) Validate() core.Error {
	if input.SujetoID == "" {
		return core.NewInvalidArgumentError("ID del sujeto es requerido")
	}
	return nil
}
