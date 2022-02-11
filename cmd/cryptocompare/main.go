package main

import (
	"errors"
	"log"
	"os"

	"github.com/vench/cryptocompare/internal/config"
	"github.com/vench/cryptocompare/internal/logger"
	"github.com/vench/cryptocompare/internal/storage/mysql"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		if errors.Is(err, config.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("failed to read app config: %w", err)
	}

	ll, err := logger.New()
	if err != nil {
		log.Fatalf("failed to create logger: %w", err)
	}
	defer ll.Sync()

	storage, err := mysql.New(&appConfig.Mysql)
	if err != nil {
		log.Fatalf("failed to create storage: %w", err)
	}
	defer storage.Close()

	log.Println(appConfig.ExternalServices.UrlCryptoCompare)
}
