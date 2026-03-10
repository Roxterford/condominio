package usecase

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type Handler[ctx context.Context, input any, output any] interface {
	Exec(ctx ctx, input input) (output, core.Error)
}

type WithInput[input any] Handler[context.Context, input, any]

type WithContextInput[ctx context.Context, input any] Handler[ctx, input, any]

type WithOutput[output any] Handler[context.Context, any, output]

type WithContextOutput[ctx context.Context, output any] Handler[ctx, any, output]

type WithInOut[input any, output any] Handler[context.Context, input, output]

type WithContextInOut[ctx context.Context, input any, output any] Handler[ctx, input, output]
