package entities

type Currency struct {
	FromSymbol string
	ToSymbol   string

	CHANGE24HOUR    string
	CHANGEPCT24HOUR string
	OPEN24HOUR      string
	VOLUME24HOUR    string
	VOLUME24HOURTO  string
	LOW24HOUR       string
	HIGH24HOUR      string
	PRICE           string
	SUPPLY          string
	MKTCAP          string
}
