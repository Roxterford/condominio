// TODO: Considera renombrar el paquete
package context

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/exception"

	"github.com/Sanaruca/condominio/internal/core/session"
)

// Definimos un tipo privado para las llaves.
// Nadie fuera de este paquete puede crear una llave de tipo 'ctxKey'.
type ctxKey int

const (
	// Usamos iota para asignar valores únicos internamente.
	user_key ctxKey = iota
	session_key
	correlation_key
	aggregate_key
)

var (
	ErrUnauthorized        = exception.New(exception.UNAUTHORIZED, "se requiere sesión de usuario")
	ErrSinPrivilegiosAdmin = exception.New(exception.FORBIDDEN,
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

// CorrelationIDKey es la key para acceder al correlation ID desde el context.
func CorrelationIDKey() any { return correlation_key }

// InjectCorrelationID inserta el correlation ID en el context.
func InjectCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlation_key, correlationID)
}

// CorrelationIDFromContext extrae el correlation ID del context.
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(correlation_key)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

// AggregateIDKey es la key para acceder al aggregate ID desde el context.
func AggregateIDKey() any { return aggregate_key }

// InjectAggregateID inserta el aggregate ID en el context.
func InjectAggregateID(ctx context.Context, aggregateID string) context.Context {
	return context.WithValue(ctx, aggregate_key, aggregateID)
}

// AggregateIDFromContext extrae el aggregate ID del context.
func AggregateIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(aggregate_key)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

// MustCorrelationID extrae el correlation ID o genera uno nuevo si no existe.
func MustCorrelationID(ctx context.Context) string {
	if id, ok := CorrelationIDFromContext(ctx); ok && id != "" {
		return id
	}
	// Generar UUID v7-like (timestamp + random) - usando google/uuid
	return "gen-" + uuid.New().String()
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
