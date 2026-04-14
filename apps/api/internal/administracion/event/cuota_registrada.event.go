package event

type CuotaRegistrada struct {
	ID string
}

func NewCuotaRegistrada(id string) CuotaRegistrada {
	return CuotaRegistrada{
		ID: id,
	}
}

func (c CuotaRegistrada) EventName() string {
	return "cuota.registrada"
}

func (c CuotaRegistrada) Payload() map[string]any {
	return map[string]any{
		"id": c.ID,
	}
}
