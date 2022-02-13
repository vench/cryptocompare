package server

import "github.com/vench/cryptocompare/internal/entities"

type currencyRawResponse struct {
	CHANGE24HOUR    float64 `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      float64 `json:"OPEN24HOUR"`
	VOLUME24HOUR    float64 `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  float64 `json:"VOLUME24HOURTO"`
	LOW24HOUR       float64 `json:"LOW24HOUR"`
	HIGH24HOUR      float64 `json:"HIGH24HOUR"`
	PRICE           float64 `json:"PRICE"`
	MKTCAP          float64 `json:"MKTCAP"`
	SUPPLY          float64 `json:"SUPPLY"`
}

type currencyDisplayResponse struct {
	CHANGE24HOUR    string `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR string `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      string `json:"OPEN24HOUR"`
	VOLUME24HOUR    string `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  string `json:"VOLUME24HOURTO"`
	LOW24HOUR       string `json:"LOW24HOUR"`
	HIGH24HOUR      string `json:"HIGH24HOUR"`
	PRICE           string `json:"PRICE"`
	SUPPLY          string `json:"SUPPLY"`
	MKTCAP          string `json:"MKTCAP"`
}

type priceResponse struct {
	Raw     map[string]map[string]*currencyRawResponse     `json:"RAW,omitempty"`
	Display map[string]map[string]*currencyDisplayResponse `json:"DISPLAY,omitempty"`
}

func MakeCurrencyResponse(result []*entities.Currency) interface{} {
	response := &priceResponse{
		Raw:     make(map[string]map[string]*currencyRawResponse),
		Display: make(map[string]map[string]*currencyDisplayResponse),
	}

	for _, item := range result {
		// raw data
		fm, ok := response.Raw[item.FromSymbol]
		if !ok {
			fm = make(map[string]*currencyRawResponse)
		}

		fm[item.ToSymbol] = &currencyRawResponse{
			CHANGE24HOUR:    item.CHANGE24HOUR,
			CHANGEPCT24HOUR: item.CHANGEPCT24HOUR,
			OPEN24HOUR:      item.OPEN24HOUR,
			VOLUME24HOUR:    item.VOLUME24HOUR,
			VOLUME24HOURTO:  item.VOLUME24HOURTO,
			LOW24HOUR:       item.LOW24HOUR,
			HIGH24HOUR:      item.HIGH24HOUR,
			PRICE:           item.PRICE,
			SUPPLY:          item.SUPPLY,
			MKTCAP:          item.MKTCAP,
		}

		response.Raw[item.FromSymbol] = fm

		// display data
		dm, ok := response.Display[item.FromSymbol]
		if !ok {
			dm = make(map[string]*currencyDisplayResponse)
		}

		dm[item.ToSymbol] = &currencyDisplayResponse{
			CHANGE24HOUR:    MoneyDollarFormat(item.CHANGE24HOUR),
			CHANGEPCT24HOUR: MoneyFormat(item.CHANGEPCT24HOUR, "", 2),
			OPEN24HOUR:      MoneyDollarFormat(item.OPEN24HOUR),
			VOLUME24HOUR:    MoneyBitcoinFormat(item.VOLUME24HOUR),
			VOLUME24HOURTO:  MoneyDollarFormat(item.VOLUME24HOURTO),
			LOW24HOUR:       MoneyDollarFormat(item.LOW24HOUR),
			HIGH24HOUR:      MoneyDollarFormat(item.HIGH24HOUR),
			PRICE:           MoneyDollarFormat(item.PRICE),
			SUPPLY:          MoneyBitcoinFormat(item.SUPPLY),
			MKTCAP:          MoneyDollarFormat(item.MKTCAP),
		}

		response.Display[item.FromSymbol] = dm
	}

	return response
}
