package model

import (
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
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

func ProveedorFromDomain(p proveedor.Proveedor) *Proveedor {
	email := p.Email().String()
	telefono := p.Telefono().String()

	return &Proveedor{
		ID:            p.ID(),
		Rif:           p.Rif().String(),
		Nombre:        p.Nombre(),
		Email:         &email,
		Telefono:      &telefono,
		Direccion:     p.Direccion(),
		CreadoEn:      p.CreadoEn(),
		ActualizadoEn: p.ActualizadoEn(),
	}
}
