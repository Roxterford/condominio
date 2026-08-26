package command

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

type RegistrarUnidadDTO struct {
	Codigo           string
	Estado           estadounidad.EstadoDeUnidad
	TitularPrimarioID *string
	ContactoID       *string
	Descripcion      string
}

type RegistrarUnidad usecase.Handler[cc.BaseContext, RegistrarUnidadDTO, *unidad.Unidad]

type registrarUnidad struct {
	unidadRepo   unidad.UnidadRepository
	unidadFactory *unidad.UnidadFactory
	sujetoRepo   sujeto.SujetoRepository
}

func NewRegistrarUnidad(
	unidadRepository unidad.UnidadRepository,
	unidadFactory *unidad.UnidadFactory,
	sujetoRepository sujeto.SujetoRepository,
) RegistrarUnidad {

	if unidadRepository == nil {
		panic("unidadRepository is nil")
	}

	if unidadFactory == nil {
		panic("unidadFactory is nil")
	}

	if sujetoRepository == nil {
		panic("sujetoRepository is nil")
	}

	return &registrarUnidad{
		unidadRepo:   unidadRepository,
		unidadFactory: unidadFactory,
		sujetoRepo:   sujetoRepository,
	}
}

func (uc *registrarUnidad) Exec(
	ctx cc.BaseContext,
	input RegistrarUnidadDTO,
) (*unidad.Unidad, core.Error) {

	codigo := unidad.UnidadCodigo(input.Codigo)
	if codigo == "" {
		return nil, core.NewValidationError("el codigo es requerido")
	}

	existe, err := uc.unidadRepo.ExistsCodigo(ctx, codigo)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, core.NewValidationError("ya existe una unidad con este codigo")
	}

	var titular sujeto.Titular
	if input.TitularPrimarioID != nil {
		s, e := uc.sujetoRepo.ObtenerPorID(ctx, sujeto.SujetoID(*input.TitularPrimarioID))
		if e != nil {
			return nil, e
		}
		if s == nil {
			return nil, core.NewValidationError("el titular primario no existe")
		}
		titular = s.AsTitular()
	}

	var contacto *sujeto.Persona
	if input.ContactoID != nil {
		s, e := uc.sujetoRepo.ObtenerPorID(ctx, sujeto.SujetoID(*input.ContactoID))
		if e != nil {
			return nil, e
		}
		if s == nil {
			return nil, core.NewValidationError("el contacto no existe")
		}
		contacto = s.AsPersona()
	}

	nueva, ferr := uc.unidadFactory.Nueva(codigo, input.Estado, titular, contacto, input.Descripcion)
	if ferr != nil {
		return nil, ferr
	}

	if gerr := uc.unidadRepo.Guardar(ctx, nueva); gerr != nil {
		return nil, gerr
	}

	return nueva, nil
}

func (input *RegistrarUnidadDTO) Validate() core.Error {
	if input.Codigo == "" {
		return core.NewValidationError("el codigo es requerido")
	}
	return nil
}
