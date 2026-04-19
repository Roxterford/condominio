package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota/estadoproyecto"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodecuota"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

type Recaudacion struct {
	Cuota              string
	Tipo               tipodecuota.TipoDeCuota
	Monto              int
	Mes                mes.Mes
	Anio               int
	Unidades           int
	UnidadesAplicadas  int `gorm:"column:unidades_aplicadas"`
	UnidadesSolventes  int `gorm:"column:unidades_solventes"`
	UnidadesPendientes int `gorm:"column:unidades_pendientes"`
	TotalEstimado      int
	Recaudado          int
	Pendiente          int
	PagosAsociados     int `gorm:"column:pagos_asociados"`
}

func (t Recaudacion) TableName() string {
	return "recaudacion"
}

func (r Recaudacion) ToDomainRecaudacion(qf *quantity.QuantityFactory) *cuota.Recaudacion {
	return &cuota.Recaudacion{
		Cuota:              cuota.CuotaID(r.Cuota),
		Tipo:               r.Tipo,
		Mes:                r.Mes,
		Anio:               r.Anio,
		MontoCuota:         qf.Assemble(int64(r.Monto)),
		MontoRecaudado:     qf.Assemble(int64(r.Recaudado)),
		MontoPendiente:     qf.Assemble(int64(r.Pendiente)),
		MontoEstimado:      qf.Assemble(int64(r.TotalEstimado)),
		Unidades:           r.Unidades,
		UnidadesAplicadas:  r.UnidadesAplicadas,
		UnidadesSolventes:  r.UnidadesSolventes,
		UnidadesPendientes: r.UnidadesPendientes,
		PagosAsociados:     r.PagosAsociados,
	}
}

type Proyecto struct {
	Titulo         string
	Cuota          string
	Estado         estadoproyecto.EstadoDeProyecto
	Descripcion    string
	Justificacion  string
	FechaLimite    time.Time
	InteresPorMora int `gorm:"column:interes_por_mora"`
	Registro       time.Time
	Actualizacion  time.Time
	RegistradoPor  string `gorm:"column:registrado_por"`
	ActualizadoPor string `gorm:"column:actualizado_por"`
}

func (p Proyecto) ToDomainProyecto(factory *cuota.ProyectoFactory) cuota.Proyecto {
	return *factory.Assemble(
		p.Titulo,
		p.Descripcion,
		p.Justificacion,
		p.Estado,
		p.FechaLimite,
		int64(p.InteresPorMora),
		p.Registro,
		p.Actualizacion,
		p.RegistradoPor,
		p.ActualizadoPor,
	)
}

type Cuota struct {
	ID             string
	Tipo           tipodecuota.TipoDeCuota
	Monto          int
	Mes            int
	Anio           int
	Registro       time.Time
	RegistradoPor  string `gorm:"column:registrado_por"`
	Actualizacion  time.Time
	ActualizadoPor string `gorm:"column:actualizado_por"`

	// Relations
	Proyecto *Proyecto `gorm:"foreignKey:cuota"`
}

func (t Cuota) TableName() string {
	return "cuotas"
}

func (c Cuota) ToDomainCuota(factory *cuota.CuotaFactory) cuota.Cuota {

	switch c.Tipo {
	case tipodecuota.Regular:
		return factory.AssembleRegular(
			c.ID,
			c.Monto,
			c.Mes,
			c.Anio,
			c.Registro,
			c.Actualizacion,
			c.RegistradoPor,
		)

	case tipodecuota.Especial:
		return factory.AssembleEspecial(
			c.ID,
			c.Mes,
			c.Anio,
			c.Monto,
			c.Proyecto.Titulo,
			c.Proyecto.Descripcion,
			c.Proyecto.Justificacion,
			c.Proyecto.Estado.String(),
			c.Proyecto.FechaLimite,
			int64(c.Proyecto.InteresPorMora),
			c.Proyecto.Registro,
			c.Proyecto.Actualizacion,
			c.Proyecto.RegistradoPor,
			c.Proyecto.ActualizadoPor,
		)
	}

	return nil
}

type Proveedor struct {
	ID            string
	Rif           string
	Nombre        string
	Email         string
	Telefono      string
	Direccion     *string
	Registro      time.Time
	Actualizacion time.Time
}

func (t Proveedor) TableName() string {
	return "proveedores"
}

func (t Proveedor) ToDomainProveedor(factory *proveedor.ProveedorFactory) *proveedor.Proveedor {
	return factory.Assemble(
		t.ID,
		t.Rif,
		t.Nombre,
		t.Email,
		t.Telefono,
		t.Direccion,
		t.Registro,
		t.Actualizacion,
	)

}

type Deuda struct {
	ID     string
	Unidad string
	Cuota  string
	// Este valor es el monto de la cuota
	Monto int
	// Este valor es el monto restante de la cuota e ira reduciendose a medida que
	// se realicen los pagos
	Deuda         int
	Estado        string
	Registro      time.Time
	Actualizacion time.Time
}

type DestinoDePago struct {
	ID        string
	Pago      string
	Deuda     string
	Destinado int
	Fecha     time.Time
}

type Gasto struct {
	ID             string
	Concepto       string
	Proveedor      string
	Cuota          *string
	Monto          int
	Moneda         moneda.Moneda
	Tasa           int
	Fecha          time.Time
	Descripcion    *string
	Registro       time.Time
	Registrado_por string
}

func (t Gasto) ToDomainGasto(factory *gasto.GastoFactory) *gasto.Gasto {

	return factory.Assemble(
		t.ID,
		t.Concepto,
		t.Proveedor,
		t.Cuota,
		t.Monto,
		t.Moneda,
		t.Tasa,
		t.Fecha,
		t.Descripcion,
		t.Registrado_por,
		t.Registro,
	)

}

func (Gasto) TableName() string {
	return "internal_gastos"
}

type GastoView struct {
	Gasto
	Total int
}

func (t GastoView) TableName() string {
	return "gastos"
}

type UnidadesTotales struct {
	TotalUnidades         int
	UnidadesActivas       int
	UnidadesInhabitadas   int
	UnidadesExentas       int
	UnidadesEnLitigio     int
	UnidadesSuspendidas   int
	UnidadesPreventa      int
	UnidadesConPendientes int
	UnidadesSolventes     int
	TotalPendiente        int
	TotalAsignado         int
}

func (t UnidadesTotales) TableName() string {
	return "unidades_totales"
}
