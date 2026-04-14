package estadovilla

type EstadoDeVilla string

const (
	// Activa: Es el estado por defecto. Según el Artículo 12 de la LPH, todo propietario
	// debe contribuir a los gastos comunes en proporción a su alícuota.
	Activa EstadoDeVilla = "ACTIVA"

	// Inhabitada: Genera deuda según la LPH, pero permite al sistema lógica para
	// omitir gastos variables (consumos) si la asamblea o el documento lo permiten.
	Inhabitada EstadoDeVilla = "INHABITADA"

	// Exenta: No genera deuda. Se usa para bienes comunes (casa conserje, áreas administrativas)
	// o por acuerdo unánime en el Documento de Condominio.
	Exenta EstadoDeVilla = "EXENTA"

	// EnLitigio: La villa sigue generando deuda (alícuota), pero este estado puede
	// bloquear la emisión de recibos estándar o estados de cuenta simples, ya que
	// el proceso de cobro está en manos de abogados (Cobro Extrajudicial o Judicial).
	EnLitigio EstadoDeVilla = "EN_LITIGIO"

	// Suspendida: Se usa para villas que, por razones legales o de fuerza mayor
	// (ej. una intervención de protección civil o colapso estructural), no pueden
	// ser sujetas a cobro temporalmente mientras se resuelve su estatus.
	Suspendida EstadoDeVilla = "SUSPENDIDA"

	// Preventa: La villa está registrada en el sistema pero aún no ha sido entregada
	// al propietario final. Los gastos suelen cargarse a la constructora o quedar
	// en suspenso hasta la protocolización.
	Preventa EstadoDeVilla = "PREVENTA"
)
