package scheduler

import (
	"context"
	"time"

	"github.com/vench/cryptocompare/internal/storage"

	"go.uber.org/zap"

	"github.com/vench/cryptocompare/internal/config"
)

type Scheduler struct {
	conf           *config.AppConfig
	logger         *zap.Logger
	currencyReader storage.CurrencyReader
	currencyWriter storage.CurrencyWriter

	done chan struct{}
}

// NewScheduler create new instance Scheduler.
func NewScheduler(
	logger *zap.Logger,
	conf *config.AppConfig,
	currencyReader storage.CurrencyReader,
	currencyWriter storage.CurrencyWriter) (*Scheduler, error) {
	return &Scheduler{
		conf:           conf,
		logger:         logger,
		currencyReader: currencyReader,
		currencyWriter: currencyWriter,
		done:           make(chan struct{}),
	}, nil
}

func (s *Scheduler) Run(ctx context.Context) error {
	t := time.NewTicker(s.conf.Scheduler.TickInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := s.parseApiCryptocompare(); err != nil {
				s.logger.Error("failed to parse api cryptocompare", zap.Error(err))
			}
		case <-ctx.Done():
			return nil
		case <-s.done:
			return nil
		}
	}

	return nil
}

func (s *Scheduler) Close() error {
	close(s.done)
	return nil
}
