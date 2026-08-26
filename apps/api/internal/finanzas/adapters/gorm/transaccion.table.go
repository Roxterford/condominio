package gorm

import (
	"time"
)

// Transaccion es la proyeccion GORM del contenedor que agrupa dos o mas
// operaciones relacionadas. Solo existe para casos excepcionales (compensacion).
type Transaccion struct {
	ID            string `gorm:"primaryKey"`
	Fecha         time.Time
	Concepto      string
	RegistradoPor string
	Registro      time.Time
}

func (Transaccion) TableName() string {
	return "transacciones"
}

// TransaccionOperacion es la tabla puente 1:N entre el contenedor y sus
// operaciones. Cada operacion pertenece como maximo a una transaccion.
type TransaccionOperacion struct {
	ID            string `gorm:"primaryKey"`
	TransaccionID string
	OperacionID   string
	Posicion      int
}

func (TransaccionOperacion) TableName() string {
	return "transaccion_operaciones"
}
