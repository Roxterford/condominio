package unidad_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

func TestFacturacionPolicyPorDefectoGeneraDeuda(t *testing.T) {
	policy := unidad.NewFacturacionPolicyPorDefecto()

	tests := []struct {
		estado         estadounidad.EstadoDeUnidad
		generaRegular  bool
		generaEspecial bool
	}{
		{estadounidad.Activa, true, true},
		{estadounidad.Inhabitada, true, true},
		{estadounidad.EnLitigio, true, true},
		{estadounidad.Exenta, false, false},
		{estadounidad.Suspendida, false, false},
		{estadounidad.Preventa, false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.estado), func(t *testing.T) {
			assert.Equal(t, tt.generaRegular, policy.GeneraDeudaPorCuotaRegular(tt.estado))
			assert.Equal(t, tt.generaEspecial, policy.GeneraDeudaPorCuotaEspecial(tt.estado))
		})
	}
}
