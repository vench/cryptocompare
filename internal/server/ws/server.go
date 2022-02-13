package ws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vench/cryptocompare/internal/config"
	"github.com/vench/cryptocompare/internal/storage"
	"go.uber.org/zap"
)

// Server contains and produce maintenance ws service.
type Server struct {
	logger  *zap.Logger
	conf    *config.AppConfig
	storage storage.CurrencyReader
}

// NewServer create instance of Server.
func NewServer(logger *zap.Logger, conf *config.AppConfig, currencyReader storage.CurrencyReader) (*Server, error) {
	return &Server{
		logger:  logger,
		conf:    conf,
		storage: currencyReader,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.conf.WS.Port),
		Handler: s.router(),
	}

	errCh := make(chan error)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	s.logger.Info("WS server is running",
		zap.Int("port", s.conf.WS.Port),
	)

	select {
	case <-ctx.Done():
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}
