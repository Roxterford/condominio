package audit

import "time"

// CreationAudit representa quién y cuándo dio origen a un registro.
// Propósito: Entidades inmutables (Pagos, Movimientos, Logs).
type CreationAudit[ID comparable] struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy ID        `json:"created_by"` // ID del Usuario o Sistema
}

// NewCreationAudit es un constructor para asegurar la validez del objeto.
func NewCreationAudit[ID comparable](userID ID) CreationAudit[ID] {
	return CreationAudit[ID]{
		CreatedAt: time.Now().UTC(),
		CreatedBy: userID,
	}
}

// FullAudit agrupa el rastro completo de una entidad mutable.
// Propósito: Entidades maestras o de configuración.
type FullAudit[ID comparable] struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy ID        `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy ID        `json:"updated_by"`
}

// Update actualiza los campos de modificación de forma atómica.
func (a *FullAudit[ID]) Update(userID ID) {
	a.UpdatedAt = time.Now().UTC()
	a.UpdatedBy = userID
}

// VersionedAudit añade un contador de cambios para evitar colisiones.
// Propósito: Sistemas con alta concurrencia o auditoría de cambios frecuentes.
type VersionedAudit[ID comparable] struct {
	FullAudit[ID]
	Version int `json:"version"` // Incrementa con cada Update
}

func (v *VersionedAudit[ID]) Increment(userID ID) {
	v.Update(userID)
	v.Version++
}
