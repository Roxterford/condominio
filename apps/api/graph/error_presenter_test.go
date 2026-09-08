package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/Sanaruca/condominio/internal/core/exception"
)

func TestErrorPresenter_NumberOverflow(t *testing.T) {
	_, overflowErr := graphql.UnmarshalInt32(int64(12312312312312300))

	ctx := graphql.WithPathContext(context.Background(), graphql.NewPathWithField("monto"))

	gqlErr := ErrorPresenter(ctx, overflowErr)

	if gqlErr.Message != "el valor del campo `monto` excede el máximo permitido" {
		t.Fatalf("mensaje inesperado: %q", gqlErr.Message)
	}

	if got, _ := gqlErr.Extensions["code"].(exception.CoreErrorCode); got != exception.INVALID_ARGUMENT {
		t.Fatalf("code inesperado: %v", gqlErr.Extensions["code"])
	}
}

func TestErrorPresenter_InternalIsMappedThrough(t *testing.T) {
	coreErr := exception.New(exception.INTERNAL, "un mensaje")
	gqlErr := ErrorPresenter(context.Background(), coreErr)

	if gqlErr == nil {
		t.Fatal("esperaba gqlerror.Error")
	}
}

func TestErrorPresenter_RawWrapsToInternal(t *testing.T) {
	gqlErr := ErrorPresenter(context.Background(), errors.New("boom"))

	if gqlErr == nil {
		t.Fatal("esperaba gqlerror.Error")
	}

	if gqlErr.Message != "[INTERNAL]: Ocurrió un error inesperado" {
		t.Fatalf("mensaje inesperado: %q", gqlErr.Message)
	}
}

func TestErrorPresenter_WrappedGqlErrorUnwrapsAndWraps(t *testing.T) {
	wrapped := &gqlerror.Error{
		Err: errors.New("boom"),
	}

	gqlErr := ErrorPresenter(context.Background(), wrapped)

	if gqlErr == nil {
		t.Fatal("esperaba gqlerror.Error")
	}

	if gqlErr.Message != "[INTERNAL]: Ocurrió un error inesperado" {
		t.Fatalf("mensaje inesperado: %q", gqlErr.Message)
	}
}

func TestErrorPresenter_EngineErrorPassesThrough(t *testing.T) {
	engineErr := &gqlerror.Error{
		Message: "Cannot query field `foo` on type `Query`",
	}

	gqlErr := ErrorPresenter(context.Background(), engineErr)

	if gqlErr == nil {
		t.Fatal("esperaba gqlerror.Error")
	}

	if gqlErr.Message != engineErr.Message {
		t.Fatalf("debería pasar tal cual, se obtuvo: %q", gqlErr.Message)
	}
}
