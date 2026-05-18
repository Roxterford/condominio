package model

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

func (input *PagoFilter) ToFilter() filter.Filter[pago.Pago] {
	return applyFilter[pago.Pago](input)
}

func PagoFromDomain(p pago.Pago) *Pago {

	return &Pago{
		ID:             p.ID(),
		Unidad:         p.Unidad(),
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

func DeudaFromDomain(d deuda.Deuda) *Deuda {

	abonos := make([]*Abono, len(d.Abonos()))

	for i, a := range d.Abonos() {
		abonos[i] = &Abono{
			Pago:  a.PagoID(),
			Monto: a.Monto().Float(),
			Fecha: a.Fecha(),
		}
	}

	return &Deuda{
		ID:       d.ID(),
		Cuota:    d.CuotaID().String(),
		Unidad:   d.Unidad(),
		Monto:    d.Monto().Float(),
		Abonos:   abonos,
		Registro: d.Registro(),
		Estado:   d.Estado(),
		Deuda:    d.Deuda().Float(),
	}
}

func UnidadesTotalesFromDomain(estadisticas unidad.Estadisticas) *UnidadesTotales {

	return &UnidadesTotales{
		TotalUnidades:         int32(estadisticas.TotalUnidades),
		UnidadesActivas:       int32(estadisticas.UnidadesActivas),
		UnidadesInhabitadas:   int32(estadisticas.UnidadesInhabitadas),
		UnidadesExentas:       int32(estadisticas.UnidadesExentas),
		UnidadesEnLitigio:     int32(estadisticas.UnidadesEnLitigio),
		UnidadesSuspendidas:   int32(estadisticas.UnidadesSuspendidas),
		UnidadesPreventa:      int32(estadisticas.UnidadesPreventa),
		UnidadesConPendientes: int32(estadisticas.UnidadesConPendientes),
		UnidadesSolventes:     int32(estadisticas.UnidadesSolventes),
		TotalPendiente:        estadisticas.TotalPendiente.Float(),
		TotalAsignado:         estadisticas.TotalAsignado.Float(),
	}

}

func UnidadFromDomain(unidad unidad.Unidad) *Unidad {

	contacto := PersonaFromDomain(unidad.Contacto())
	titular := TitularFromDomain(unidad.TitularPrimario())

	return &Unidad{
		ID:              unidad.ID(),
		Codigo:          unidad.Codigo(),
		Estado:          unidad.Estado(),
		TitularPrimario: titular,
		Contacto:        contacto,
		Deuda:           unidad.Deuda().Float(),
	}
}

func TitularFromDomain(titular sujeto.Titular) Titular {
	if titular == nil {
		return nil
	}

	ente := titular.AsEnte()
	persona := titular.AsPersona()

	if ente != nil {
		return EnteFromDomain(ente)
	}

	if persona != nil {
		return PersonaFromDomain(persona)
	}

	return nil

}

func SujetoFromDomain(sujeto sujeto.Sujeto) Sujeto {
	if sujeto == nil {
		return nil
	}

	ente := sujeto.AsEnte()
	persona := sujeto.AsPersona()

	if ente != nil {
		return EnteFromDomain(ente)
	}

	if persona != nil {
		return PersonaFromDomain(persona)
	}

	return nil

}

func EnteFromDomain(ente *sujeto.Ente) *Ente {
	if ente == nil {
		return nil
	}

	return &Ente{
		ID:            ente.ID().String(),
		RazonSocial:   ente.RazonSocial(),
		Email:         ente.Email().String(),
		Telefono:      ente.Telefono().String(),
		Cedula:        ente.Cedula().String(),
		Registro:      ente.Audit().CreatedAt,
		Actualizacion: ente.Audit().UpdatedAt,
	}

}

func PersonaFromDomain(persona *sujeto.Persona) *Persona {
	if persona == nil {
		return nil
	}

	return &Persona{
		ID:            persona.ID().String(),
		Nombres:       persona.Nombres(),
		Apellidos:     persona.Apellidos(),
		Email:         persona.Email().String(),
		Telefono:      persona.Telefono().String(),
		Cedula:        persona.Cedula().String(),
		Registro:      persona.Audit().CreatedAt,
		Actualizacion: persona.Audit().UpdatedAt,
	}

}

func RecaudacionFromDomain(recuadacion cuota.Recaudacion) Recaudacion {
	return Recaudacion{
		Moneda:             moneda.USD,
		MontoEstimado:      recuadacion.MontoEstimado.Float(),
		MontoRecaudado:     recuadacion.MontoRecaudado.Float(),
		MontoPendiente:     recuadacion.MontoPendiente.Float(),
		PagosAsociados:     int32(recuadacion.PagosAsociados),
		Unidades:           int32(recuadacion.Unidades),
		UnidadesAplicadas:  int32(recuadacion.UnidadesAplicadas),
		UnidadesSolventes:  int32(recuadacion.UnidadesSolventes),
		UnidadesPendientes: int32(recuadacion.UnidadesPendientes),
	}
}

func CuotaTypeFromDomain(cuota cuota.Cuota) CuotaType {
	if cuota == nil {
		return nil
	}
	cuota_regular := cuota.AsRegular()
	cuota_especial := cuota.AsEspecial()

	if cuota_regular != nil {
		return CuotaRegular{
			ID:            string(cuota_regular.ID()),
			Monto:         cuota_regular.Monto().Float(),
			Mes:           int32(cuota_regular.Mes()),
			Anio:          int32(cuota_regular.Anio()),
			Registro:      cuota_regular.Audit.CreatedAt,
			Actualizacion: cuota_regular.Audit.UpdatedAt,
		}
	}

	if cuota_especial != nil {
		return CuotaEspecial{
			ID:            string(cuota_especial.ID()),
			Monto:         cuota_especial.Monto().Float(),
			Mes:           int32(cuota_especial.Mes()),
			Anio:          int32(cuota_especial.Anio()),
			Registro:      cuota_especial.Audit.CreatedAt,
			Actualizacion: cuota_especial.Audit.UpdatedAt,
			Detalles: &Proyecto{
				Titulo:         cuota_especial.Detalles.Titulo(),
				Estado:         cuota_especial.Detalles.Estado(),
				Descripcion:    cuota_especial.Detalles.Descripcion(),
				Justificacion:  cuota_especial.Detalles.Justificacion(),
				FechaLimite:    cuota_especial.Detalles.FechaLimite(),
				InteresPorMora: float64(cuota_especial.Detalles.InteresPorMora().Value()),
				Registro:       cuota_especial.Detalles.Audit.CreatedAt,
				Actualizacion:  cuota_especial.Detalles.Audit.UpdatedAt,
			},
		}
	}

	return nil
}

// applyFilter es una función privada que convierte un input a Filter usando genéricos
func applyFilter[T filter.Filterable](input any) filter.Filter[T] {
	nill := *filter.NewFilter[T](nil)
	if input == nil {
		return nill
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nill
	}

	var inputMap map[string]any
	if err := json.Unmarshal(jsonBytes, &inputMap); err != nil {
		return nill
	}

	return *filter.NewFilter[T](inputMap)
}

func (input *UnidadFilter) ToFilter() filter.Filter[unidad.Unidad] {
	return applyFilter[unidad.Unidad](input)
}

func (input *CuotaFilter) ToFilter() filter.Filter[cuota.CuotaBase] {
	return applyFilter[cuota.CuotaBase](input)
}

func (input *ObtenerProveedoresDto) ToFilter() *filter.Filter[proveedor.Proveedor] {
	if input == nil {
		return nil
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil
	}

	var inputMap map[string]any
	if err := json.Unmarshal(jsonBytes, &inputMap); err != nil {
		return nil
	}

	return filter.NewFilter[proveedor.Proveedor](inputMap)
}

func GastoYProveedorFromDomain(
	gasto gasto.Gasto,
	proveedor proveedor.Proveedor,
) *GastoWithProveedor {
	email := proveedor.Email().String()
	telefono := proveedor.Telefono().String()
	direccion := proveedor.Direccion()

	return &GastoWithProveedor{
		ID:       string(gasto.ID()),
		Concepto: gasto.Concepto(),
		Proveedor: &Proveedor{
			ID:            proveedor.ID(),
			Rif:           proveedor.Rif().String(),
			Nombre:        proveedor.Nombre(),
			Email:         &email,
			Telefono:      &telefono,
			Direccion:     direccion,
			CreadoEn:      proveedor.CreadoEn(),
			ActualizadoEn: proveedor.ActualizadoEn(),
		},
		Cuota:         gasto.Cuota(),
		Monto:         (gasto.Monto()).Float(),
		Moneda:        gasto.Moneda(),
		Tasa:          (gasto.Tasa()).Float(),
		Total:         (gasto.Total()).Float(),
		Fecha:         gasto.Fecha(),
		Descripcion:   gasto.Descripcion(),
		Registro:      gasto.Audit().CreatedAt,
		RegistradoPor: gasto.Audit().CreatedBy,
	}
}

func GastoFromDomain(gasto gasto.Gasto) *Gasto {
	return &Gasto{
		ID:        string(gasto.ID()),
		Concepto:  gasto.Concepto(),
		Proveedor: gasto.Proveedor(),
		Cuota:     gasto.Cuota(),
		// Monto:         int32(gasto.Monto()),
		Moneda: string(gasto.Moneda()),
		// Tasa:          int32(gasto.Tasa()),
		// Total:         int32(gasto.Total()),
		Fecha:         gasto.Fecha(),
		Descripcion:   gasto.Descripcion(),
		Registro:      gasto.Audit().CreatedAt,
		RegistradoPor: gasto.Audit().CreatedBy,
	}
}

func (p *Paginator) ToDomainPaginator() common.Paginator {
	if p == nil {
		return common.Paginator{}
	}
	return common.Paginator{
		Page:  int(p.Page),
		Limit: int(p.Limit),
	}
}
