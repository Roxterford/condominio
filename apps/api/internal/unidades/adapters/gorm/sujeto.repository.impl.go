package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

type sujetoRepository struct {
	db      *gorm.DB
	factory *sujeto.SujetoFactory
}

func NewSujetoRepository(db *gorm.DB, factory *sujeto.SujetoFactory) sujeto.SujetoRepository {
	if factory == nil {
		panic("factory is nil")
	}

	return &sujetoRepository{db, factory}
}

// ObtenerPorID implements [sujeto.SujetoRepository].
func (r *sujetoRepository) ObtenerPorID(
	ctx context.Context,
	id sujeto.SujetoID,
) (sujeto.Sujeto, core.Error) {

	s, err := gorm.G[Sujeto](r.db).Where("id = ?", id).Preload("Sujeto", nil).Take(ctx)

	if exception.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	switch s.Tipo {
	case PERSONA_NATURAL:

		persona := s.ToDoaminPersona(r.factory)

		return &persona, nil

	case ENTE_JURIDICO:
		var razon_social string

		if s.RazonSocial != nil {
			razon_social = *s.RazonSocial
		}

		return r.factory.AssembleEnte(
			s.ID,
			s.DocumentoIdentidad,
			razon_social,
			s.DocumentoIdentidad,
			s.Email,
			s.Representante.ToDoaminPersona(r.factory),
		), nil
	default:
		return nil, core.NewError(exception.CONFLICT, "tipo de sujeto desconocido")
	}

}
