package common

import (
	"regexp"
	"strings"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/exception"
)

var (
	// Error para validación de email inválido
	ErrInvalidEmail = core.NewInvalidArgumentError(
		"El email proporcionado no es válido. Formato esperado: nombre@dominio.com",
	)

	ErrForbiddenDomain = exception.New(
		exception.FORBIDDEN,
		"El dominio proporcionado no está permitido.",
	)
)

// Regex estándar para validación
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type EmailFactory struct {
	allowed_domains map[string]bool
}

type Email struct {
	address string
}

func (e Email) Address() string {
	return e.address
}

func (e Email) String() string { return e.address }

func (e Email) Validate() core.Error {

	if !emailRegex.MatchString(e.address) {
		return ErrInvalidEmail
	}

	return nil
}

// Equals permite comparar dos emails por valor, no por referencia
func (e Email) Equals(other Email) bool {
	return e.address == other.address
}

func NewEmailFactory(domains []string) *EmailFactory {
	allowed := make(map[string]bool)
	for _, d := range domains {
		allowed[strings.ToLower(d)] = true
	}
	return &EmailFactory{allowed_domains: allowed}
}

func (f *EmailFactory) New(v string) (Email, core.Error) {
	email := strings.ToLower(strings.TrimSpace(v))

	// 1. Validación robusta con Regex
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	// 2. Extraer dominio (ahora es seguro porque el regex garantiza que hay un '@')
	parts := strings.Split(email, "@")
	domain := parts[len(parts)-1]

	// 3. Validar contra la Lista Blanca
	if len(f.allowed_domains) > 0 && !f.allowed_domains[domain] {
		return Email{}, ErrForbiddenDomain
	}

	return Email{address: email}, nil
}

// Assemble ignores policy, only reconstructs
func (f *EmailFactory) Assemble(raw string) Email {

	return Email{address: raw}
}
