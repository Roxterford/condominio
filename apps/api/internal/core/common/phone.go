package common

import (
	"fmt"
	"strings"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

const (
	DEFAULT_PHONE_SUBSCRIBER_LENGTH = 6
)

var (
	ErrInvalidPhoneFormat = errors.New(errors.INVALID_ARGUMENT, "formato de teléfono inválido")
	ErrCountryNotAllowed  = errors.New(errors.INVALID_ARGUMENT, "país no permitido")
	ErrProviderNotAllowed = errors.New(errors.INVALID_ARGUMENT, "proveedor no permitido")
)

// PhoneFactory es una fabrica que valida y construye objetos Phone
type PhoneFactory struct {
	allowed_countries map[string]bool
	allowed_providers map[string]bool // Opcional
}

func NewPhoneFactory(countries []string, providers []string) *PhoneFactory {
	cMap := make(map[string]bool)
	for _, c := range countries {
		cMap[c] = true
	}
	pMap := make(map[string]bool)
	for _, p := range providers {
		pMap[p] = true
	}
	return &PhoneFactory{
		allowed_countries: cMap,
		allowed_providers: pMap,
	}
}

func (f *PhoneFactory) New(raw string) (Phone, core.Error) {
	clean := strings.ReplaceAll(raw, " ", "")
	parts := strings.Split(clean, "-")

	if len(parts) != 3 {
		return Phone{}, ErrInvalidPhoneFormat
	}

	cc, prov, sub := parts[0], parts[1], parts[2]

	// Validar país
	if len(f.allowed_countries) > 0 && !f.allowed_countries[cc] {
		return Phone{}, ErrCountryNotAllowed
	}

	// Validar proveedor (opcional)
	if len(f.allowed_providers) > 0 && !f.allowed_providers[prov] {
		return Phone{}, ErrProviderNotAllowed
	}

	// Validación de longitud mínima básica para el suscriptor
	if len(sub) < DEFAULT_PHONE_SUBSCRIBER_LENGTH {
		return Phone{}, ErrInvalidPhoneFormat
	}

	return Phone{
		country_code: cc,
		provider:     prov,
		subscriber:   sub,
	}, nil
}

type Phone struct {
	country_code string // ej: "58"
	provider     string // ej: "414"
	subscriber   string // ej: "1234567"
}

func (p Phone) CountryCode() string { return p.country_code }
func (p Phone) Provider() string    { return p.provider }
func (p Phone) Subscriber() string  { return p.subscriber }

// FullNumber devuelve el número completo con el signo +
func (p Phone) FullNumber() string {
	return fmt.Sprintf("+%s%s%s", p.country_code, p.provider, p.subscriber)
}

func (p Phone) String() string {
	return p.FullNumber()
}

func (f *PhoneFactory) Assemble(raw string) (Phone, core.Error) {
	clean := strings.ReplaceAll(raw, " ", "")
	parts := strings.Split(clean, "-")

	if len(parts) != 3 {
		return Phone{}, ErrInvalidPhoneFormat
	}

	cc, prov, sub := parts[0], parts[1], parts[2]

	return Phone{
		country_code: cc,
		provider:     prov,
		subscriber:   sub,
	}, nil
}
