package model

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

func CuotaTypeFromDomain(cuota cuota.Cuota) CuotaType {
	if cuota == nil {
		return nil
	}
	cuota_regular := cuota.AsRegular()
	cuota_especial := cuota.AsEspecial()

	if cuota_regular != nil {
		return CuotaRegular{
			ID:            string(cuota_regular.ID),
			Monto:         int32(cuota_regular.Monto),
			Mes:           int32(cuota_regular.Mes),
			Anio:          int32(cuota_regular.Anio),
			Registro:      cuota_regular.Audit.CreatedAt,
			Actualizacion: cuota_regular.Audit.UpdatedAt,
			Pagos:         &PagosCount{},
			Villas:        &VillasCount{},
		}
	}

	if cuota_especial != nil {
		return CuotaEspecial{
			ID:            string(cuota_especial.ID),
			Monto:         int32(cuota_especial.Monto),
			Mes:           int32(cuota_especial.Mes),
			Anio:          int32(cuota_especial.Anio),
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
			Pagos:  &PagosCount{},
			Villas: &VillasCount{},
		}
	}

	return nil
}

func (input *CuotaFilter) ToFilter() filter.Filter[cuota.CuotaBase] {
	nill := *filter.NewFilter[cuota.CuotaBase](nil)
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

	return *filter.NewFilter[cuota.CuotaBase](inputMap)
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
		Monto:         int32(gasto.Monto()),
		Moneda:        gasto.Moneda(),
		Tasa:          int32(gasto.Tasa()),
		Total:         int32(gasto.Total()),
		Fecha:         gasto.Fecha(),
		Descripcion:   gasto.Descripcion(),
		Registro:      gasto.Audit().CreatedAt,
		RegistradoPor: gasto.Audit().CreatedBy,
	}
}

func GastoFromDomain(gasto gasto.Gasto) *Gasto {
	return &Gasto{
		ID:            string(gasto.ID()),
		Concepto:      gasto.Concepto(),
		Proveedor:     gasto.Proveedor(),
		Cuota:         gasto.Cuota(),
		Monto:         int32(gasto.Monto()),
		Moneda:        string(gasto.Moneda()),
		Tasa:          int32(gasto.Tasa()),
		Total:         int32(gasto.Total()),
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
