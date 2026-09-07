package loaders

import (
	"context"
	"fmt"

	"github.com/Sanaruca/condominio/graph/model"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
	database "github.com/Sanaruca/condominio/internal/shared/adapters/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type unidadReader struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
}

func GetUnidad(ctx context.Context, unidadID string) (*model.Unidad, error) {
	loaders := For(ctx)
	return loaders.Unidad.Load(ctx, unidadID)
}

func GetUnidades(ctx context.Context, unidadIDs []string) ([]*model.Unidad, error) {
	loaders := For(ctx)
	return loaders.Unidad.LoadAll(ctx, unidadIDs)
}

func GetUnidadID(ctx context.Context, unidadID string) (*model.UnidadIdentifiers, error) {
	loaders := For(ctx)
	return loaders.UnidadIdentifiers.Load(ctx, unidadID)
}

func GetUnidadIDs(ctx context.Context, unidadIDs []string) ([]*model.UnidadIdentifiers, error) {
	loaders := For(ctx)
	return loaders.UnidadIdentifiers.LoadAll(ctx, unidadIDs)
}

func (r *unidadReader) getIdentifiers(
	ctx context.Context,
	ids []string,
) ([]*model.UnidadIdentifiers, []error) {

	rows, err := gorm.G[database.UnidadInfo](r.db).
		Where(fmt.Sprintf("%s.id IN ?", new(database.UnidadInfo).TableName()), ids).
		Or(fmt.Sprintf("%s.codigo IN ?", new(database.UnidadInfo).TableName()), ids).
		// TODO: ordenar para que no pase 1, 10, 2, 20 ...
		Order(fmt.Sprintf("%s.codigo asc", new(database.UnidadInfo).TableName())).
		Find(ctx)

	if err != nil {
		return nil, []error{err}
	}

	identifiers := make([]*model.UnidadIdentifiers, len(rows))

	for i, unidad := range rows {

		identifiers[i] = &model.UnidadIdentifiers{
			ID:     unidad.ID,
			Codigo: unidad.Codigo,
		}
	}

	return identifiers, nil
}

func (r *unidadReader) getUnidades(ctx context.Context, ids []string) ([]*model.Unidad, []error) {

	rows, err := gorm.G[database.UnidadInfo](r.db).
		Where(fmt.Sprintf("%s.id IN ?", new(database.UnidadInfo).TableName()), ids).
		Or(fmt.Sprintf("%s.codigo IN ?", new(database.UnidadInfo).TableName())).
		Joins(clause.LeftJoin.Association("Contacto"), nil).
		Joins(clause.LeftJoin.Association("TitularPrimario"), nil).
		// TODO: ordenar para que no pase 1, 10, 2, 20 ...
		Order(fmt.Sprintf("%s.codigo asc", new(database.UnidadInfo).TableName())).
		Find(ctx)

	if err != nil {
		return nil, []error{err}
	}

	unidades := make([]*model.Unidad, len(rows))

	for i, unidad := range rows {

		var contacto *model.Persona
		var titular model.Titular

		if unidad.Contacto != nil {

			var c = unidad.Contacto

			nombres := "NULL"
			apellidos := ""

			if c.Nombres != nil && c.Apellidos != nil {
				nombres = *c.Nombres
				apellidos = *c.Apellidos
			} else {
				logger.ErrorCtx(
					ctx,
					nil,
					"Nombre del contacto de la unidad desconocido",
					"unidad.contacto.nombres", c.Nombres,
					"unidad.contacto.apellidos", c.Apellidos,
				)
			}

			contacto = &model.Persona{
				ID:          c.ID,
				Nombres:     nombres,
				Apellidos:   apellidos,
				Email:       c.Email,
				Telefono:    c.Telefono,
				Cedula:      c.DocumentoIdentidad,
				Registro:    c.Registro,
				DisplayName: c.DisplayName(),
			}
		}

		if unidad.TitularPrimario != nil {

			t := unidad.TitularPrimario

			switch t.Tipo {
			case database.TipoDeSujetoPersonaNatural:

				nombres := "NULL"
				apellidos := ""

				if t.Nombres != nil && t.Apellidos != nil {
					nombres = *t.Nombres
					apellidos = *t.Apellidos
				} else {
					logger.ErrorCtx(
						ctx,
						nil,
						"Nombre del tiular de la unidad desconocido",
						"unidad.tiular.nombres", t.Nombres,
						"unidad.tiular.apellidos", t.Apellidos,
					)
				}

				titular = &model.Persona{
					ID:          t.ID,
					Nombres:     nombres,
					Apellidos:   apellidos,
					Email:       t.Email,
					Telefono:    t.Telefono,
					Cedula:      t.DocumentoIdentidad,
					Registro:    t.Registro,
					DisplayName: t.DisplayName(),
				}

			case database.TipoDeSujetoEnteJuridico:

				razon := "NULL"

				if t.RazonSocial != nil {
					razon = *t.RazonSocial
				} else {
					logger.ErrorCtx(ctx, nil, "Razon social del titular de la unidad desconocido", "unidad.titular.razon_social", t.RazonSocial)
				}

				titular = &model.Ente{
					ID:            t.ID,
					RazonSocial:   razon,
					Email:         t.Email,
					Telefono:      t.Telefono,
					Cedula:        t.DocumentoIdentidad,
					Representante: contacto,
					Registro:      t.Registro,
					DisplayName:   t.DisplayName(),
				}
			default:
				logger.ErrorCtx(
					ctx,
					nil,
					"No se pudo determinar el tipo de titular de la unidad",
					"unidad.titular.tipo",
					unidad.TitularPrimario.Tipo,
				)
			}

		}

		unidades[i] = &model.Unidad{
			ID:              unidad.ID,
			Codigo:          unidad.Codigo,
			Estado:          unidad.Estado,
			TitularPrimario: titular,
			Titulares:       []model.Titular{},
			Contacto:        contacto,
			Deuda:           r.qf.Assemble(int64(unidad.DeudaTotal)).Float(),
			Wallet:          r.qf.Assemble(int64(unidad.Wallet)).Float(),
		}
	}

	return unidades, nil
}
