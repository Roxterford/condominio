package villas

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type VillaRepository interface {
	Exists(ctx context.Context, villa int) (bool, core.Error)
}
