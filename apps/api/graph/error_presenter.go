package graph

import (
	"context"
	"errors"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

// ErrorPresenter es el safety net de errores no manejados. Se ejecuta para
// cada error GraphQL antes de serializarlo.
//
//  1. Los errores crudos (SQL, redis, red, etc.) que escapan sin wrapear se
//     envuelven con exception.Wrap, mapeando errores conocidos via el registry
//     y cayendo en INTERNAL con mensaje genérico si no se conocen.
//  2. Los CoreError INTERNAL (no mapeados) se loggean automáticamente con su
//     causa original; el resto de códigos (validation, not found, etc.) son
//     esperados y no saturan el log.
//  3. Los errores de unmarshal de gqlgen (p.ej. desbordamiento de int32) se
//     traducen a mensajes amigables.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	// Si ya es un CoreError, no se re-envuelve: solo se loggea si es INTERNAL.
	var coreErr exception.CoreError
	if errors.As(err, &coreErr) {
		if coreErr.Code() == exception.INTERNAL {
			logger.ErrorCtx(ctx, coreErr.Cause(), "error no mapeado (internal)",
				"message", coreErr.Message(),
				"operation", operationName(ctx),
				"field", fieldPath(ctx),
			)
		}

		return graphql.DefaultErrorPresenter(ctx, err)
	}

	// Errores de entrada (unmarshal de gqlgen): p.ej. valores que desbordan
	// int32. Se traducen a un mensaje amigable en lugar de exponer el detalle
	// técnico al cliente.
	var overflowErr *graphql.NumberOverflowError
	if errors.As(err, &overflowErr) {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)

		logger.ErrorCtx(ctx, err, "valor fuera de rango de entrada",
			"message", gqlErr.Message,
			"path", pathString(gqlErr.Path),
		)

		msg := "el valor enviado excede el máximo permitido"
		if field := lastPathName(gqlErr.Path); field != "" {
			msg = "el valor del campo `" + field + "` excede el máximo permitido"
		}

		return &gqlerror.Error{
			Err:     err,
			Message: msg,
			Path:    gqlErr.Path,
			Extensions: map[string]interface{}{
				"code": exception.INVALID_ARGUMENT,
			},
		}
	}

	// Errores del engine de gqlgen (parse de query, validación de variables,
	// etc.) llegan como *gqlerror.Error con Err == nil y ya son presentables.
	// Se pasan tal cual.
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) && gqlErr.Err == nil {
		return graphql.DefaultErrorPresenter(ctx, err)
	}

	// Error crudo sin envolver (librería externa: gorm, redis, http, etc.).
	// Se envuelve ahora: el registry mapea errores conocidos (gorm.ErrRecordNotFound,
	// redis.Nil, etc.) y el resto cae en INTERNAL con mensaje genérico.
	raw := unwrapError(err)
	wrapped := exception.Wrap(raw)

	return graphql.DefaultErrorPresenter(ctx, wrapped)
}

// unwrapError extrae el error original de un *gqlerror.Error (que es como
// gqlgen envuelve los errores de los resolvers antes de llamar al presenter).
func unwrapError(err error) error {
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) && gqlErr.Err != nil {
		return gqlErr.Err
	}
	return err
}

func operationName(ctx context.Context) string {
	if graphql.HasOperationContext(ctx) {
		return graphql.GetOperationContext(ctx).OperationName
	}
	return ""
}

func fieldPath(ctx context.Context) string {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return ""
	}
	return fc.Object + "." + fc.Field.Name
}

func lastPathName(p ast.Path) string {
	if len(p) == 0 {
		return ""
	}
	if n, ok := p[len(p)-1].(ast.PathName); ok {
		return string(n)
	}
	return ""
}

func pathString(p ast.Path) string {
	parts := make([]string, len(p))
	for i, el := range p {
		if n, ok := el.(ast.PathName); ok {
			parts[i] = string(n)
		} else {
			parts[i] = "? "
		}
	}
	return strings.Join(parts, ".")
}
