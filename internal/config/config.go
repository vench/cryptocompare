package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/jessevdk/go-flags"
)

var ErrHelp = errors.New("help")

type AppConfig struct {
	Name string `long:"name" description:"Application name" default:"App name"`

	HTTP HTTPServer `group:"HTTP server configuration" namespace:"http"`

	Mysql Mysql `group:"Mysql configuration" namespace:"mysql"`

	Scheduler Scheduler `group:"Scheduler configuration" namespace:"scheduler"`

	CryptoCompare CryptoCompare `group:"External services cryptocompare" namespace:"crypto_compare"`
}

type Mysql struct {
	ConnectionString string `long:"connection_string" description:"String connection MYSQL" default:"root:admin@tcp(127.0.0.1:3306)/test"`
}

type HTTPServer struct {
	Port int `long:"port" description:"Port HTTP server" default:"8090"`
}

type CryptoCompare struct {
	FromSymbols []string `long:"from_symbols" description:"From symbols" default:"BTC" default:"XRP" default:"ETH" default:"BCH" default:"EOS" default:"LTC" default:"XMR" default:"DASH"`
	ToSymbols   []string `long:"to_symbols" description:"To symbols" default:"USD" default:"EUR" default:"GBP" default:"JPY" default:"RUR"`

	Url string `long:"url" description:"Address of cryptocompare.com api" default:"https://min-api.cryptocompare.com/data/pricemultifull"`
}

type Scheduler struct {
	TickInterval time.Duration `long:"tick_interval" description:"Tick interval scheduler" default:"5m"`
}

func NewAppConfig() (*AppConfig, error) {
	var config AppConfig
	if _, err := flags.Parse(&config); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return nil, ErrHelp
		}
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
