package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/villas"
	"gorm.io/gorm"
)

type VillaGORMRepository struct {
	db *gorm.DB
}

func NewVillaGORMRepository(db *gorm.DB) villas.VillaRepository {
	return &VillaGORMRepository{
		db: db,
	}
}

// Exists implements [villas.VillaRepository].
func (r *VillaGORMRepository) Exists(ctx context.Context, villa int) (bool, core.Error) {

	_, err := gorm.G[Villa](r.db).Select("id").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, core.WrapError(err)
	}

	return true, nil
}
