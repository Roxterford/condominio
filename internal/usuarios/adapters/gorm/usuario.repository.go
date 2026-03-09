package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/usuarios"
	"gorm.io/gorm"
)

type UsuarioGORMRepository struct {
	db      *gorm.DB
	factory usuarios.UsuarioFactory
}

func NewUsuarioGORMRepository(
	db *gorm.DB,
	factory usuarios.UsuarioFactory,
) usuarios.UsuarioRepository {
	return &UsuarioGORMRepository{db: db, factory: factory}
}

// GetByEmail implements [usuarios.UsuarioRepository].
func (r *UsuarioGORMRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*usuarios.Usuario, core.Error) {
	usuario, err := gorm.G[Usuario](r.db).Where("email = ?", email).Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	return r.factory.Assemble(
		usuario.ID,
		usuario.Email,
		usuario.Password,
		usuario.Nombre,
		usuario.Apellido,
	)

}

// Guardar implements [usuarios.UsuarioRepository].
func (r *UsuarioGORMRepository) Guardar(ctx context.Context, usuario *usuarios.Usuario) core.Error {

	existe, err := r.usuarioExiste(ctx, usuario.ID())

	if err != nil {
		return err
	}

	if existe {
		return r.actualizarUsurioExistente(ctx, usuario)
	}

	return r.insertarNuevoUsurio(ctx, usuario)

}

func (r *UsuarioGORMRepository) actualizarUsurioExistente(
	ctx context.Context,
	usuario *usuarios.Usuario,
) core.Error {

	_, err := gorm.G[Usuario](
		r.db,
	).Where("id = ?", usuario.ID()).
		Updates(ctx, *mapToUsuarioTable(usuario))

	return core.WrapError(err)

}
func (r *UsuarioGORMRepository) insertarNuevoUsurio(
	ctx context.Context,
	usuario *usuarios.Usuario,
) core.Error {

	err := gorm.G[Usuario](r.db).Create(ctx, mapToUsuarioTable(usuario))

	return core.WrapError(err)

}

func (r *UsuarioGORMRepository) usuarioExiste(
	ctx context.Context,
	usuarioID string,
) (bool, core.Error) {

	_, err := gorm.G[Usuario](r.db).Select("id").Where("id = ?", usuarioID).Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, core.WrapError(err)
	}

	return true, nil

}
