// Deprecated: Mappers legacy de pagos.
package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
)

func (input *PagoFilter) ToFilter() filter.Filter[pago.Pago] {
	return applyFilter[pago.Pago](input)
}

func PagoFromDomain(p pago.Pago) *Pago {

	return &Pago{
		ID:             p.ID(),
		Unidad:         string(p.Unidad()),
		Fecha:          p.Fecha(),
		Metodo:         p.Metodo(),
		Monto:          p.Monto().Float(),
		Referencia:     p.Referencia(),
		Moneda:         p.Moneda(),
		Tasa:           p.Tasa().Float(),
		Total:          p.Total().Float(),
		Registro:       p.Firma().RegistradoEn(),
		RegistradoPor:  p.Firma().RegistradoPor(),
		Actualizacion:  p.Firma().ActualizadoEn(),
		ActualizadoPor: p.Firma().ActualizadoPor(),
		Destinado:      p.SaldoDestinado().Float(),
		Disponible:     p.SaldoDisponible().Float(),
	}

}
