package query_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/stretchr/testify/assert"
)

func TestObtenerProveedor(t *testing.T) {

	db := utils.SetupInMemoryTestDB(t)

	if err := db.Create(&administracion.Proveedor{
		ID:     "1",
		Nombre: "Proveedor 1",
	}).Error; err != nil {
		t.Fatal(err)
	}

	uc := query.NewObtenerProveedor()

	proveedor, err := uc.Exec(
		context.New(t.Context(), db),
		query.ObtenerProveedorDTO{ProveedorID: "1"},
	)

	assert.NoError(t, err)
	assert.NotNil(t, proveedor)
	assert.Equal(t, "1", proveedor.ID)
	assert.Equal(t, "Proveedor 1", proveedor.Nombre)

}
