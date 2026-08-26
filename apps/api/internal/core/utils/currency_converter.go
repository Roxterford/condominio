package utils

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type CurrencyConverter struct {
	qf *quantity.QuantityFactory
}

func NewCurrencyConverter(quantityFactory *quantity.QuantityFactory) *CurrencyConverter {

	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}

	return &CurrencyConverter{
		qf: quantityFactory,
	}

}

func (c CurrencyConverter) Convert(value, change int) quantity.Quantity {

	_value := c.qf.Assemble(int64(value))
	_change := c.qf.Assemble(int64(change))

	return _value.HappyDiv(_change)
}
