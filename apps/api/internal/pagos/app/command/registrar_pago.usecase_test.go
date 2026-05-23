package command_test

import (
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegistrarPagoDTO(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	ref := "ABC123"

	tests := []struct {
		name    string
		dto     command.RegistrarPagoDTO
		wantErr bool
	}{
		{
			name: "válido con todos los campos",
			dto: command.RegistrarPagoDTO{
				Unidad:     unidad.UnidadCodigo("1"),
				Fecha:      &now,
				Metodo:     "EFECTIVO", // asume que existe un método válido
				Referencia: &ref,
				Monto:      100,
				Tasa:       10,
				Moneda:     "USD", // asume que existe una moneda válida
			},
			wantErr: false,
		},
		{
			name: "Unidad inválida (0)",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("0"), Fecha: &now, Metodo: metododepago.Efectivo,
				Monto: 100, Tasa: 10, Moneda: moneda.USD,
			},
			wantErr: true,
		},
		{
			name: "Monto inválido (0)",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("1"), Fecha: &now, Metodo: metododepago.Efectivo,
				Monto: 0, Tasa: 10, Moneda: moneda.VED,
			},
			wantErr: true,
		},
		{
			name: "Tasa inválida (0)",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("1"), Fecha: &now, Metodo: metododepago.Efectivo,
				Monto: 100, Tasa: 0, Moneda: moneda.VED,
			},
			wantErr: true,
		},
		{
			name: "Referencia demasiado larga",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("1"), Fecha: &now, Metodo: metododepago.Efectivo,
				Monto: 100, Tasa: 10, Moneda: moneda.VED,
				Referencia: func() *string {
					s := string(make([]byte, 60)) // 60 chars
					return &s
				}(),
			},
			wantErr: true,
		},
		{
			name: "Fecha futura inválida",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("1"), Fecha: &future, Metodo: metododepago.Efectivo,
				Monto: 100, Tasa: 10, Moneda: moneda.VED,
			},
			wantErr: true,
		},
		{
			name: "Fecha nil (se asigna automáticamente)",
			dto: command.RegistrarPagoDTO{
				Unidad: unidad.UnidadCodigo("1"), Fecha: nil, Metodo: metododepago.Efectivo,
				Monto: 100, Tasa: 10, Moneda: moneda.VED,
			},
			wantErr: false,
		},
	}

	for _, tst := range tests {
		test := tst
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.dto.Validate()
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
