package common

const (
	DEFAULT_PAGE_SIZE = 20
)

type Paginator struct {
	Page  int
	Limit int
}

type Paginated[T any] struct {
	// Registros Paginados
	Data []T
	// Total de registros
	Total int
	// Pagina actual
	Page int
	// Paginas totales
	Pages int
	// Limite de registros en la paginacion
	Limit int
}

func NewPaginated[T any](data []T, total int, pagninator Paginator) *Paginated[T] {
	pagninator.Sanitize()
	pages := CalculatePages(total, pagninator.Limit)

	return &Paginated[T]{
		Data:  data,
		Total: total,
		Page:  pagninator.Page,
		Pages: pages,
		Limit: pagninator.Limit,
	}

}

func (p Paginator) Offset() int {
	p.Sanitize()
	return (p.Page - 1) * p.Limit
}

func (p *Paginator) Sanitize() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = DEFAULT_PAGE_SIZE
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

func CalculatePages(total, limit int) int {
	if total <= 0 || limit <= 0 {
		return 1
	}

	pages := (total + limit - 1) / limit

	return pages

}
