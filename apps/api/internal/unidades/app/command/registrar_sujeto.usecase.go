package command

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

type RegistrarSujetoDTO struct {
	Tipo              sujeto.TipoDeSujeto
	DocumentoIdentidad string
	Nombres           string
	Apellidos         string
	RazonSocial       string
	RepresentanteID   *string
	Email             string
	Telefono          string
}

type RegistrarSujeto usecase.Handler[cc.BaseContext, RegistrarSujetoDTO, sujeto.Sujeto]

type registrarSujeto struct {
	sujetoRepo    sujeto.SujetoRepository
	sujetoFactory *sujeto.SujetoFactory
}

func NewRegistrarSujeto(
	sujetoRepository sujeto.SujetoRepository,
	sujetoFactory *sujeto.SujetoFactory,
) RegistrarSujeto {

	if sujetoRepository == nil {
		panic("sujetoRepository is nil")
	}

	if sujetoFactory == nil {
		panic("sujetoFactory is nil")
	}

	return &registrarSujeto{
		sujetoRepo:    sujetoRepository,
		sujetoFactory: sujetoFactory,
	}
}

func (uc *registrarSujeto) Exec(
	ctx cc.BaseContext,
	input RegistrarSujetoDTO,
) (sujeto.Sujeto, core.Error) {

	existe, err := uc.sujetoRepo.ExistsDocumento(ctx, input.DocumentoIdentidad)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, core.NewValidationError("ya existe un sujeto con este documento de identidad")
	}

	existe, err = uc.sujetoRepo.ExistsEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, core.NewValidationError("ya existe un sujeto con este email")
	}

	registrador := ""
	if ctx.Session() != nil && ctx.Session().Usuario() != nil {
		registrador = ctx.Session().Usuario().ID
	}

	representanteID := ""
	if input.RepresentanteID != nil {
		representanteID = *input.RepresentanteID
	}

	nuevo, ferr := uc.sujetoFactory.Nuevo(
		input.Tipo,
		input.DocumentoIdentidad,
		input.Nombres,
		input.Apellidos,
		input.RazonSocial,
		representanteID,
		input.Email,
		input.Telefono,
		registrador,
	)
	if ferr != nil {
		return nil, ferr
	}

	if gerr := uc.sujetoRepo.Guardar(ctx, nuevo); gerr != nil {
		return nil, gerr
	}

	return nuevo, nil
}

func (input *RegistrarSujetoDTO) Validate() core.Error {
	if input.DocumentoIdentidad == "" {
		return core.NewValidationError("el documento de identidad es requerido")
	}
	if input.Email == "" {
		return core.NewValidationError("el email es requerido")
	}
	if input.Telefono == "" {
		return core.NewValidationError("el telefono es requerido")
	}
	return nil
}
