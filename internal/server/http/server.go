package http

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/vench/cryptocompare/internal/config"
	"go.uber.org/zap"
)

// Server contains and produce maintance web service.
type Server struct {
	logger *zap.Logger
	conf   *config.AppConfig
}

// NewServer create instance of Server.
func NewServer(logger *zap.Logger, conf *config.AppConfig) (*Server, error) {
	return &Server{
		logger: logger,
		conf:   conf,
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

func userValueUint64(rCtx *fasthttp.RequestCtx, key string) (uint64, bool) {
	value := rCtx.UserValue(key)
	sid, ok := value.(string)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
