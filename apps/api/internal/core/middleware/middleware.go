package middleware

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type Middleware[Ctx context.Context, Req any, Res any] struct {
	handler func(ctx Ctx, request Req, next usecase.Handler[Ctx, Req, Res]) (Res, core.Error)
	next    usecase.Handler[Ctx, Req, Res]
}

func New[Ctx context.Context, Req any, Res any](
	handler func(ctx Ctx, request Req, next usecase.Handler[Ctx, Req, Res]) (Res, core.Error),
	next usecase.Handler[Ctx, Req, Res],
) Middleware[Ctx, Req, Res] {

	return Middleware[Ctx, Req, Res]{
		handler: handler,
		next:    next,
	}
}

func (md Middleware[Ctx, Req, Res]) Exec(
	ctx Ctx,
	req Req,
) (res Res, err core.Error) {
	return md.handler(ctx, req, md.next)
}

func (md Middleware[Ctx, Req, Res]) GetHandler() func(ctx Ctx, request Req, next usecase.Handler[Ctx, Req, Res]) (Res, core.Error) {
	return md.handler
}

func (md *Middleware[Ctx, Req, Res]) SetNext(next usecase.Handler[Ctx, Req, Res]) {
	md.next = next
}

func UseMiddleware[Ctx context.Context, Req any, Res any](
	usecase usecase.Handler[Ctx, Req, Res],
	middleware ...func(ctx Ctx, request Req, next usecase.Handler[Ctx, Req, Res]) (Res, core.Error),
) (res usecase.Handler[Ctx, Req, Res]) {

	if len(middleware) < 1 {
		return usecase
	}

	res = New(
		middleware[len(middleware)-1],
		usecase,
	)
	// Iterate through the middleware in reverse order
	for i := len(middleware) - 2; i >= 0; i-- {
		// Create a new handler that wraps the previous one
		res = New(
			middleware[i],
			res,
		)
	}

	return res
}
