package model

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

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
