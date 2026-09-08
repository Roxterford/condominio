package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
)

// ErrorPresenter es el safety net de errores no manejados. Se ejecuta para
// cada error GraphQL antes de serializarlo. Los errores INTERNAL (no mapeados
// o inesperados) se loggean automáticamente con su causa original; el resto de
// códigos (validation, not found, etc.) son esperados y no saturan el log.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	if coreErr, ok := err.(exception.CoreError); ok && coreErr.Code() == exception.INTERNAL {
		logger.ErrorCtx(ctx, coreErr.Cause(), "error no mapeado (internal)",
			"message", coreErr.Message(),
			"operation", operationName(ctx),
			"field", fieldPath(ctx),
		)
	}

	return graphql.DefaultErrorPresenter(ctx, err)
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
