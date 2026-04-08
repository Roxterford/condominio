package cuota

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota/estadoproyecto"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
)

type Proyecto struct {
	estado           estadoproyecto.EstadoDeProyecto
	descripcion      string
	justificacion    string
	fecha_limite     time.Time
	interes_por_mora common.Percentage
	Audit            audit.FullAudit[string]
}

func (p *Proyecto) Estado() estadoproyecto.EstadoDeProyecto { return p.estado }
func (p *Proyecto) Descripcion() string                     { return p.descripcion }
func (p *Proyecto) Justificacion() string                   { return p.justificacion }
func (p *Proyecto) FechaLimite() time.Time                  { return p.fecha_limite }
func (p *Proyecto) InteresPorMora() common.Percentage {
	return p.interes_por_mora
}
func (p *Proyecto) Registro() time.Time      { return p.Audit.CreatedAt }
func (p *Proyecto) Actualizacion() time.Time { return p.Audit.UpdatedAt }
