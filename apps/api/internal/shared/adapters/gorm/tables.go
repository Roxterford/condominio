package gorm

import (
	"time"

	estadoproyecto "github.com/Sanaruca/condominio/internal/administracion/models/cuota/estadoproyecto"
	estadodeuda "github.com/Sanaruca/condominio/internal/administracion/models/deuda/estadodeuda"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodecuota"
	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestinoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipoperacion"
	estadounidad "github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

// TipoDeSujeto solo existe localmente en el adapter de Sujeto; no tiene
// representacion de dominio reutilizable.
type TipoDeSujeto string

const (
	TipoDeSujetoPersonaNatural TipoDeSujeto = "PERSONA_NATURAL"
	TipoDeSujetoEnteJuridico   TipoDeSujeto = "ENTE_JURIDICO"
)

// Usuario -> usuarios
type Usuario struct {
	ID       string `gorm:"primaryKey"`
	Email    string
	Password string
}

func (Usuario) TableName() string { return "usuarios" }

// Proyecto -> proyectos
type Proyecto struct {
	Titulo         string
	Cuota          *string `gorm:"primaryKey;column:cuota"`
	Estado         estadoproyecto.EstadoDeProyecto
	Descripcion    string
	Justificacion  string
	FechaLimite    time.Time
	InteresPorMora int
	Registro       time.Time
	RegistradoPor  string `gorm:"column:registrado_por"`
	Actualizacion  time.Time
	ActualizadoPor string `gorm:"column:actualizado_por"`
}

func (Proyecto) TableName() string { return "proyectos" }

// Unidad -> unidades
type Unidad struct {
	ID              string `gorm:"primaryKey"`
	Codigo          string
	Estado          estadounidad.EstadoDeUnidad
	TitularPrimario *string `gorm:"column:titular_primario"`
	Contacto        *string `gorm:"column:contacto"`
	Descripcion     *string
}

func (Unidad) TableName() string { return "unidades" }

// Sujeto -> sujetos
type Sujeto struct {
	ID                 string `gorm:"primaryKey"`
	Tipo               TipoDeSujeto
	DocumentoIdentidad string
	Nombres            *string
	Apellidos          *string
	RazonSocial        *string
	Representante      *string `gorm:"column:representante"`
	Email              string
	Telefono           string
	Registro           time.Time
}

func (Sujeto) TableName() string { return "sujetos" }

// Titularidad -> titularidades
type Titularidad struct {
	ID      string `gorm:"primaryKey"`
	Titular string `gorm:"column:titular"`
	Unidad  string `gorm:"column:unidad"`
}

func (Titularidad) TableName() string { return "titularidades" }

// IOperacion -> operaciones
type IOperacion struct {
	ID            string `gorm:"primaryKey"`
	Fecha         time.Time
	Concepto      string
	Monto         int
	Moneda        moneda.Moneda
	Metodo        metodoperacion.MetodoDeOperacion
	Tasa          int
	Tipo          tipoperacion.TipoDeOperacion
	Rol           roldestinoperacion.RolDestinoDeOperacion
	Cuota         *string `gorm:"column:cuota"`
	UnidadCodigo  *string `gorm:"column:unidad_codigo"`
	Proveedor     *string `gorm:"column:proveedor"`
	RegistradoPor string  `gorm:"column:registrado_por"`
	Registro      time.Time
}

func (IOperacion) TableName() string { return "operaciones" }

// ITransaccion -> transacciones
type ITransaccion struct {
	ID            string `gorm:"primaryKey"`
	Fecha         time.Time
	Concepto      string
	RegistradoPor string `gorm:"column:registrado_por"`
	Registro      time.Time
}

func (ITransaccion) TableName() string { return "transacciones" }

// ITransaccionOperacion -> transaccion_operaciones
type ITransaccionOperacion struct {
	ID            string `gorm:"primaryKey"`
	TransaccionID string `gorm:"column:transaccion_id"`
	OperacionID   string `gorm:"column:operacion_id"`
	Posicion      int
}

func (ITransaccionOperacion) TableName() string { return "transaccion_operaciones" }

// Gasto (view) -> gastos
type Gasto struct {
	Operacion     string  `gorm:"primaryKey"`
	Transaccion   *string `gorm:"column:transaccion"`
	Fecha         time.Time
	Concepto      string
	Monto         int
	Moneda        moneda.Moneda
	Metodo        metodoperacion.MetodoDeOperacion
	Tasa          int
	Tipo          tipoperacion.TipoDeOperacion
	Rol           roldestinoperacion.RolDestinoDeOperacion
	Cuota         *string `gorm:"column:cuota"`
	UnidadCodigo  *string `gorm:"column:unidad_codigo"`
	ProveedorID   *string `gorm:"column:proveedor_id"`
	RegistradoPor string  `gorm:"column:registrado_por"`
	Registro      time.Time
}

func (Gasto) TableName() string { return "gastos" }

// DestinoDePago -> destino_de_pagos
type DestinoDePago struct {
	ID        string `gorm:"primaryKey"`
	Operacion string `gorm:"column:operacion"`
	Deuda     string `gorm:"column:deuda"`
	Destinado int
	Fecha     time.Time
}

func (DestinoDePago) TableName() string { return "destino_de_pagos" }

// Cuota -> cuotas
type Cuota struct {
	ID             string `gorm:"primaryKey"`
	Tipo           tipodecuota.TipoDeCuota
	Monto          int
	Mes            int
	Anio           int
	Registro       time.Time
	RegistradoPor  string `gorm:"column:registrado_por"`
	Actualizacion  time.Time
	ActualizadoPor string `gorm:"column:actualizado_por"`
}

func (Cuota) TableName() string { return "cuotas" }

// Proveedor -> proveedores
type Proveedor struct {
	ID            string `gorm:"primaryKey"`
	RIF           string
	Nombre        string
	Email         string
	Telefono      string
	Direccion     *string
	Registro      time.Time
	Actualizacion time.Time
}

func (Proveedor) TableName() string { return "proveedores" }

// IDeuda -> internal_deudas
type IDeuda struct {
	ID            string `gorm:"primaryKey"`
	Unidad        string `gorm:"column:unidad"`
	Cuota         string `gorm:"column:cuota"`
	Registro      time.Time
	Actualizacion time.Time
}

func (IDeuda) TableName() string { return "internal_deudas" }

// Deuda (view) -> deudas
type Deuda struct {
	ID            string `gorm:"primaryKey"`
	UnidadID      string `gorm:"column:unidad_id"`
	UnidadCodigo  string `gorm:"column:unidad_codigo"`
	Cuota         string `gorm:"column:cuota"`
	Monto         int
	Deuda         int
	Estado        estadodeuda.EstadoDeDeuda
	Registro      time.Time
	Actualizacion time.Time
}

func (Deuda) TableName() string { return "deudas" }

// Recaudacion (view) -> recaudacion
type Recaudacion struct {
	Cuota              string `gorm:"primaryKey"`
	Monto              int
	Mes                int
	Anio               int
	Unidades           int
	UnidadesAplicadas  int `gorm:"column:unidades_aplicadas"`
	UnidadesSolventes  int `gorm:"column:unidades_solventes"`
	UnidadesPendientes int `gorm:"column:unidades_pendientes"`
	TotalEstimado      int `gorm:"column:total_estimado"`
	Recaudado          int
	Pendiente          int
	PagosAsociados     int `gorm:"column:pagos_asociados"`
}

func (Recaudacion) TableName() string { return "recaudacion" }

// UnidadInfo (view) -> unidades_info
type UnidadInfo struct {
	ID               string `gorm:"primaryKey"`
	Codigo           string
	Estado           estadounidad.EstadoDeUnidad
	Contacto         *string
	TitularPrimario  *string `gorm:"column:titular_primario"`
	Descripcion      *string
	DeudaTotal       int    `gorm:"column:deuda_total"`
	EstadoCuenta     string `gorm:"column:estado_cuenta"`
	CuotasPendientes int    `gorm:"column:cuotas_pendientes"`
}

func (UnidadInfo) TableName() string { return "unidades_info" }

// TasaDeCambio (view) -> tasas_de_cambio
type TasaDeCambio struct {
	Origen   string
	Tasa     int
	Fecha    time.Time
	OrigenID string `gorm:"primaryKey;column:origen_id"`
}

func (TasaDeCambio) TableName() string { return "tasas_de_cambio" }
