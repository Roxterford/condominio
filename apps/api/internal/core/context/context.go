package context

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"

	"github.com/Sanaruca/condominio/internal/core/session"
)

// Definimos un tipo privado para las llaves.
// Nadie fuera de este paquete puede crear una llave de tipo 'ctxKey'.
type ctxKey int

const (
	// Usamos iota para asignar valores únicos internamente.
	user_key ctxKey = iota
	session_key
)

var (
	ErrUnauthorized        = errors.New(errors.UNAUTHORIZED, "se requiere sesión de usuario")
	ErrSinPrivilegiosAdmin = errors.New(errors.FORBIDDEN,
		"privilegios de administrador requeridos",
	)
)

// --- Jerarquía de Contextos ---

type BaseContext interface {
	context.Context
	Session() session.Session
}

type AuthContext interface {
	context.Context
	Session() session.SessionWithUser
}

// AdminContext hereda de UserContext (ya tiene Session() -> SessionWithUser)
type AdminContext interface {
	AuthContext
}

// //////////////////

type appContext struct {
	context.Context
}

// --- Wrappers de Promoción ---

type userContextWrapper struct{ *appContext }
type adminContextWrapper struct{ *appContext }

// Implementación de Session para BaseContext
func (c *appContext) Session() session.Session {
	u, _ := c.Value(user_key).(*session.CredencialDeUsuario)
	return &sessionImpl{u}
}

// Implementación de Session para User/Admin (Garantiza Valor)
func (c *appContext) sessionWithUser() session.SessionWithUser {
	u := c.Value(user_key).(*session.CredencialDeUsuario) // El upgrade ya validó que no es nil
	return &sessionWithUserImpl{*u}
}

func (w *userContextWrapper) Session() session.SessionWithUser  { return w.sessionWithUser() }
func (w *adminContextWrapper) Session() session.SessionWithUser { return w.sessionWithUser() }

// --- Implementaciones de Session ---

type sessionImpl struct{ u *session.CredencialDeUsuario }

func (s *sessionImpl) Usuario() *session.CredencialDeUsuario { return s.u }

type sessionWithUserImpl struct{ u session.CredencialDeUsuario }

func (s *sessionWithUserImpl) Usuario() session.CredencialDeUsuario { return s.u }

///////

type FluentWrap struct {
	ctx context.Context
}

func Wrap(ctx context.Context) FluentWrap {
	return FluentWrap{ctx}
}

// InjectUser es un helper para el Middleware: inserta el payload usando la constante privada.
func InjectUser(ctx context.Context, u *session.CredencialDeUsuario) context.Context {
	return context.WithValue(ctx, user_key, u)
}

func (f FluentWrap) AsBase() (BaseContext, core.Error) {
	if f.ctx.Value(user_key) == nil {
		return nil, nil
	}
	_, ok := f.ctx.Value(user_key).(*session.CredencialDeUsuario)
	if !ok {
		return nil, ErrUnauthorized
	}

	return &appContext{f.ctx}, nil
}

func (f FluentWrap) AsAuth() (AuthContext, core.Error) {
	u, ok := f.ctx.Value(user_key).(*session.CredencialDeUsuario)
	if !ok || u == nil {
		return nil, ErrUnauthorized
	}
	return &userContextWrapper{&appContext{f.ctx}}, nil
}

// TODO: Implementar validación de admin cuando se defina la estructura de roles
func (f FluentWrap) AsAdmin() (AdminContext, core.Error) {
	return f.AsAuth()
}
