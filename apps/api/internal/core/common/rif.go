package common

import (
	"regexp"
	"strings"

	"github.com/Sanaruca/condominio/internal/core"
)

var (
	ErrRifInvalido = core.NewInvalidArgumentError("formato de RIF inválido")
)

const rifRegex = `^[VJEGP][0-9]{8,9}$`

type Rif struct {
	value string
}

func NewRif(input string) (Rif, core.Error) {
	clean := strings.ToUpper(strings.ReplaceAll(input, "-", ""))

	if !regexp.MustCompile(rifRegex).MatchString(clean) {
		return Rif{}, ErrRifInvalido
	}

	return Rif{value: clean}, nil
}

func (r Rif) String() string { return r.value }

// Métodos de conveniencia para lógica de negocio
func (r Rif) EsPersonaNatural() bool {
	return strings.HasPrefix(r.value, "V") || strings.HasPrefix(r.value, "E")
}
func (r Rif) EsEmpresa() bool { return strings.HasPrefix(r.value, "J") }

func AssembleRif(input string) Rif {
	clean := strings.ToUpper(strings.ReplaceAll(input, "-", ""))
	return Rif{value: clean}
}
