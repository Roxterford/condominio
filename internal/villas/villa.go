package villas

type Villa struct {
	id     string
	numero int
	deuda  int
}

func (v *Villa) ID() string  { return v.id }
func (v *Villa) Numero() int { return v.numero }
func (v *Villa) Deuda() int  { return v.deuda }

func (v *Villa) PoseeDeuda() bool { return v.deuda > 0 }
