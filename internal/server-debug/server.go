package serverdebug

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	server *http.Server
}

type Options struct {
	addr string
}

func NewOptions(addr string) Options {
	return Options{
		addr: addr,
	}
}

func New(opts Options) (*Server, error) {
	mux := http.NewServeMux()

	// Регистрируем pprof endpoints для отладки
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			zap.L().Error("Failed to write health response", zap.Error(err))
		}
	})

	server := &http.Server{
		Addr:              opts.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // Защита от Slowloris
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{
		logger: zap.L().Named("debug-server"),
		server: server,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting debug server", zap.String("addr", s.server.Addr))

	errChan := make(chan error, 1)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("debug server listen: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.Info("Shutting down debug server")
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		return s.server.Shutdown(shutdownCtx)
	}
}
