package entities

type Currency struct {
	FromSymbol string
	ToSymbol   string

	CHANGE24HOUR    float64
	CHANGEPCT24HOUR float64
	OPEN24HOUR      float64
	VOLUME24HOUR    float64
	VOLUME24HOURTO  float64
	LOW24HOUR       float64
	HIGH24HOUR      float64
	PRICE           float64
	MKTCAP          float64
	SUPPLY          float64
}
