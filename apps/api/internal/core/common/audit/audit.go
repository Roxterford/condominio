package audit

import "time"

// CreationAudit representa quién y cuándo dio origen a un registro.
// Propósito: Entidades inmutables (Pagos, Movimientos, Logs).
type CreationAudit[UserID comparable] struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy UserID    `json:"created_by"` // ID del Usuario o Sistema
}

// NewCreationAudit es un constructor para asegurar la validez del objeto.
func NewCreationAudit[UserID comparable](userID UserID) CreationAudit[UserID] {
	return CreationAudit[UserID]{
		CreatedAt: time.Now().UTC(),
		CreatedBy: userID,
	}
}

// FullAudit agrupa el rastro completo de una entidad mutable.
// Propósito: Entidades maestras o de configuración.
type FullAudit[UserID comparable] struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy UserID    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy UserID    `json:"updated_by"`
}

// Update actualiza los campos de modificación de forma atómica.
func (a *FullAudit[UserID]) Update(userID UserID) {
	a.UpdatedAt = time.Now().UTC()
	a.UpdatedBy = userID
}

// VersionedAudit añade un contador de cambios para evitar colisiones.
// Propósito: Sistemas con alta concurrencia o auditoría de cambios frecuentes.
type VersionedAudit[UserID comparable] struct {
	FullAudit[UserID]
	Version int `json:"version"` // Incrementa con cada Update
}

func (v *VersionedAudit[UserID]) Increment(userID UserID) {
	v.Update(userID)
	v.Version++
}
