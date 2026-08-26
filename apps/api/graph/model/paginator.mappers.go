package model

import (
	"github.com/Sanaruca/condominio/internal/core/common"
)

func (p *Paginator) ToDomainPaginator() common.Paginator {
	var paginator common.Paginator

	if p != nil {
		paginator = common.Paginator{
			Page:  int(p.Page),
			Limit: int(p.Limit),
		}
	}

	paginator.Sanitize()

	return paginator
}
