package context

import "github.com/Sanaruca/condominio/internal/usuarios"

type AdminContext struct {
	BaseContext
	Session struct {
		Usuario usuarios.UsuarioPayload
	}
}

func GetAdminContext(ctx BaseContext) AdminContext {
	usuario := ctx.session.Usuario()

	if usuario == nil {
		panic("usuario no encontrado")
	}

	return AdminContext{
		BaseContext: ctx,
		Session: struct {
			Usuario usuarios.UsuarioPayload
		}{
			Usuario: *usuario,
		},
	}
}
