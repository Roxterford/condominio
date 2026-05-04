package gorm

import (
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

// Use UnidadInfo instead for queries
type Unidad struct {
	ID                string `gorm:"primaryKey"`
	Codigo            string
	Estado            estadounidad.EstadoDeUnidad
	TitularPrimarioID *string `gorm:"column:titular_primario"`
	ContactoID        *string `gorm:"column:contacto"`
	Descripcion       *string

	TitularPrimario *Sujeto
	Contacto        *Sujeto
}

func (u *Unidad) TableName() string {
	return "unidades"
}
