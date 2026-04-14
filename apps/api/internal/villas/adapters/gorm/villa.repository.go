package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	villaPkg "github.com/Sanaruca/condominio/internal/villas/models/villa"
	"github.com/Sanaruca/condominio/internal/villas/models/villa/estadovilla"
	"gorm.io/gorm"
)

type VillaGORMRepository struct {
	db *gorm.DB
}

func NewVillaGORMRepository(db *gorm.DB) villaPkg.VillaRepository {
	return &VillaGORMRepository{
		db: db,
	}
}

// Exists implements [villas.VillaRepository].
func (r *VillaGORMRepository) Exists(ctx context.Context, villaNumero int) (bool, core.Error) {

	_, err := gorm.G[Villa](r.db).Where("numero = ?", villaNumero).Select("id").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, core.WrapError(err)
	}

	return true, nil
}

// ObtenerTodas implements [villa.VillaRepository].
func (r *VillaGORMRepository) ObtenerTodas(ctx context.Context) ([]int, core.Error) {
	var villas []Villa
	if err := r.db.WithContext(ctx).Find(&villas).Error; err != nil {
		return nil, core.WrapError(err)
	}

	numeros := make([]int, len(villas))
	for i, v := range villas {
		numeros[i] = v.Numero
	}

	return numeros, nil
}

// ObtenerEstado implements [villa.VillaRepository].
func (r *VillaGORMRepository) ObtenerEstado(ctx context.Context, villaNumero int) (estadovilla.EstadoDeVilla, core.Error) {
	var villa Villa
	if err := r.db.WithContext(ctx).Where("numero = ?", villaNumero).Take(&villa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", villaPkg.ErrVillaNoEncontrada
		}
		return "", core.WrapError(err)
	}

	return estadovilla.EstadoDeVilla(villa.Estado), nil
}
