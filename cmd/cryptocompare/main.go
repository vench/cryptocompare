package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/vench/cryptocompare/internal/storage/api"

	"github.com/vench/cryptocompare/internal/storage"

	"github.com/vench/cryptocompare/internal/service/scheduler"

	"go.uber.org/zap"

	"github.com/vench/cryptocompare/internal/server/http"

	"github.com/vench/cryptocompare/internal/config"
	"github.com/vench/cryptocompare/internal/logger"
	"github.com/vench/cryptocompare/internal/storage/mysql"

	"github.com/chapsuk/grace"
	"golang.org/x/sync/errgroup"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		if errors.Is(err, config.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("failed to read app config: %v", err)
	}

	ll, err := logger.New()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer ll.Sync() //nolint:errcheck

	storageInner, err := mysql.New(&appConfig.Mysql)
	if err != nil {
		ll.Error("failed to create mysql storage", zap.Error(err))
		return
	}
	defer storageInner.Close()

	storageOuter, err := api.New(ll, &appConfig.CryptoCompare)
	if err != nil {
		ll.Error("failed to create api storage", zap.Error(err))
		return
	}

	storageChain := storage.NewCurrencyReaderChain(storageOuter, storageInner)

	serverHTTP, err := http.NewServer(ll, appConfig, storageChain)
	if err != nil {
		ll.Error("failed to create http server", zap.Error(err))
		return
	}

	serviceScheduler, err := scheduler.NewScheduler(ll, appConfig, storageOuter, storageInner)
	if err != nil {
		ll.Error("failed to create service scheduler", zap.Error(err))
		return
	}
	defer serviceScheduler.Close()

	ctx := grace.ShutdownContext(context.Background())

	gr, appctx := errgroup.WithContext(ctx)
	gr.Go(func() error {
		return serverHTTP.Serve(appctx)
	})
	gr.Go(func() error {
		return serviceScheduler.Run(appctx)
	})

	if err := gr.Wait(); err != nil {
		ll.Error("failed to wait", zap.Error(err))
	}

	ll.Info("application has been shutdown")
}
