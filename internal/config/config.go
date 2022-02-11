package config

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
)

var ErrHelp = errors.New("help")

type AppConfig struct {
	Name string `long:"name" description:"Application name" default:"App name"`

	HTTP HTTPServer `group:"HTTP server configuration" namespace:"http"`

	Mysql Mysql `group:"Mysql configuration" namespace:"mysql"`

	ExternalServices ExternalServices `group:"External services configuration" namespace:"external_services"`
}

type Mysql struct {
	ConnectionString string `long:"connection_string" description:"String connection MYSQL" default:"root:admin@tcp(127.0.0.1:3306)/test"`
}

type HTTPServer struct {
	Port int `long:"port" description:"Port HTTP server" default:"8090"`
}

type ExternalServices struct {
	UrlCryptoCompare string `long:"url_crypto_compare" description:"Address of cryptocompare.com api" default:"https://min-api.cryptocompare.com/data/pricemultifull"`
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
