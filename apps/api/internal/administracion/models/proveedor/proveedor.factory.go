package proveedor

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/lucsky/cuid"
)

type ProveedorFactory struct {
	emailFactory *common.EmailFactory
	phoneFactory *common.PhoneFactory
}

func NewProveedorFactory(
	emailFactory *common.EmailFactory,
	phoneFactory *common.PhoneFactory,
) *ProveedorFactory {

	if emailFactory == nil {
		panic("emailFactory is nil")
	}

	if phoneFactory == nil {
		panic("phoneFactory is nil")
	}

	return &ProveedorFactory{emailFactory: emailFactory, phoneFactory: phoneFactory}
}

func (f *ProveedorFactory) Nuevo(
	rif string,
	nombre string,
	email string,
	telefono string,
	direccion *string,
) (*Proveedor, core.Error) {

	if nombre == "" {
		return nil, core.NewValidationError("el nombre no puede estar vacío")
	}

	_rif, err := common.NewRif(rif)
	if err != nil {
		return nil, err
	}

	_email, err := f.emailFactory.New(email)
	if err != nil {
		return nil, err
	}

	_telefono, err := f.phoneFactory.New(telefono)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	return &Proveedor{
		id:             cuid.New(),
		rif:            _rif,
		nombre:         nombre,
		email:          _email,
		telefono:       _telefono,
		direccion:      direccion,
		creado_en:      ahora,
		actualizado_en: ahora,
	}, nil
}

func (f *ProveedorFactory) Assemble(
	id string,
	rif string,
	nombre string,
	email string,
	telefono string,
	direccion *string,
	creado_en time.Time,
	actualizado_en time.Time,
) *Proveedor {

	_rif := common.AssembleRif(rif)
	_email, _ := f.emailFactory.Assemble(email)
	_telefono, _ := f.phoneFactory.Assemble(telefono)

	return &Proveedor{
		id:             id,
		rif:            _rif,
		nombre:         nombre,
		email:          _email,
		telefono:       _telefono,
		direccion:      direccion,
		creado_en:      creado_en,
		actualizado_en: actualizado_en,
	}
}
