package estadodeuda

type EstadoDeDeuda string

const (
	Pendiente EstadoDeDeuda = "PENDIENTE"
	Abonada   EstadoDeDeuda = "ABONADA"
	Pagada    EstadoDeDeuda = "PAGADA"
)
