package storage

import (
	"github.com/hashicorp/go-multierror"
	"github.com/vench/cryptocompare/internal/entities"
)

type CurrencyReader interface {
	GetCurrencyBy(fromSymbol, toSymbol []string) ([]*entities.Currency, error)
}

type CurrencyWriter interface {
	StoreCurrency(currency ...*entities.Currency) error
}

type CurrencyReaderChain []CurrencyReader

func NewCurrencyReaderChain(readers ...CurrencyReader) CurrencyReaderChain {
	return readers
}

func (c CurrencyReaderChain) GetCurrencyBy(fromSymbol, toSymbol []string) ([]*entities.Currency, error) {
	var errResult error
	for i := range c {
		result, err := c[i].GetCurrencyBy(fromSymbol, toSymbol)
		if err == nil {
			return result, nil
		}
		errResult = multierror.Append(errResult, err)
	}

	return nil, errResult
}
