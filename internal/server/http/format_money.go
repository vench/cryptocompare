package http

import (
	"github.com/leekchan/accounting"
)

func moneyBitcoinFormat(number float64) string {
	return moneyFormat(number, "Ƀ", 1)
}

func moneyDollarFormat(number float64) string {
	return moneyFormat(number, "$", 2)
}

const magnitudeBillion float64 = 1000000000

func moneyFormat(number float64, symbol string, precision int) string {
	ac := accounting.DefaultAccounting(symbol, precision)
	ac.SetFormat("%s %v")
	ac.SetFormatNegative("%s -%v")
	if symbol == "" {
		ac.SetFormat("%v")
		ac.SetFormatNegative("-%v")
	}

	if magnitudeBillion < number {
		// nolint:gocritic
		number = number / magnitudeBillion
		ac.SetFormat("%s %v B")
		ac.SetFormatNegative("%s -%v B")
		if symbol == "" {
			ac.SetFormat("%v B")
			ac.SetFormatNegative("-%v B")
		}
	}

	return ac.FormatMoney(number)
}
