package command_test

import (
	"fmt"
	"testing"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/core/context"
	testingutils "github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/Sanaruca/condominio/internal/villas"
	"github.com/stretchr/testify/assert"
)

func TestAplicarCuota(t *testing.T) {
	db := testingutils.SetupInMemoryTestDB(t)

	const cantidad_de_villas = 45

	villas_arr := make([]*villas.Villa, cantidad_de_villas)
	for i := range villas_arr {
		villas_arr[i] = &villas.Villa{ID: fmt.Sprintf("v%d", i+1), Numero: i + 1}
	}

	if err := db.Create(villas_arr).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Create(&administracion.Cuota{
		ID:    "c0",
		Monto: 1_00,
	}).Error; err != nil {
		t.Fatal(err)
	}

	ctx := context.New(t.Context(), db)

	uc := command.NewAplicarCuota()

	if _, err := uc.Exec(ctx, command.AplicarCuotaDTO{
		CuotaID: "c0",
	}); err != nil {
		t.Fatal(err)
	}

	var deudas_aplicadas int64
	db.Model(new(villas.IDeuda)).Count(&deudas_aplicadas).Where("id = ?", "c0")

	assert.Equal(
		t,
		int64(cantidad_de_villas),
		deudas_aplicadas,
		"La cuota no se aplico correctamente las '%d' villas",
		cantidad_de_villas,
	)
}
