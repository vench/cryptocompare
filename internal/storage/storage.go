package storage

import "github.com/vench/cryptocompare/internal/entities"

type Storage interface {
	StoreCurrency(currency *entities.Currency) error
	GetCurrencyBy(fromSymbol, toSymbol string) (*entities.Currency, error)
}
