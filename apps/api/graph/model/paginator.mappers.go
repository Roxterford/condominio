package model

import (
	"github.com/Sanaruca/condominio/internal/core/common"
)

func (p *Paginator) ToDomainPaginator() common.Paginator {
	if p == nil {
		return common.Paginator{}
	}
	return common.Paginator{
		Page:  int(p.Page),
		Limit: int(p.Limit),
	}
}
