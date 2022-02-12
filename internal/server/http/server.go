package http

import (
	"context"
	"fmt"
	"time"

	"github.com/vench/cryptocompare/internal/storage"

	"github.com/valyala/fasthttp"
	"github.com/vench/cryptocompare/internal/config"
	"go.uber.org/zap"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

// Server contains and produce maintance web service.
type Server struct {
	logger  *zap.Logger
	conf    *config.AppConfig
	storage storage.CurrencyReader
}

// NewServer create instance of Server.
func NewServer(logger *zap.Logger, conf *config.AppConfig, storage storage.CurrencyReader) (*Server, error) {
	return &Server{
		logger:  logger,
		conf:    conf,
		storage: storage,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	srv := &fasthttp.Server{
		Handler:            s.router(ctx),
		Name:               s.conf.Name + " http server",
		ReadTimeout:        time.Second,
		WriteTimeout:       time.Second,
		CloseOnShutdown:    true,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: time.Minute,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe(fmt.Sprintf(":%d", s.conf.HTTP.Port))
	}()

	s.logger.Info("HTTP server is running",
		zap.Int("port", s.conf.HTTP.Port),
	)

	select {
	case <-ctx.Done():
		if err := srv.Shutdown(); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}

	return nil
}
