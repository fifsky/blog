package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// Option is an HTTP server option.
type Option func(*Server)

// Handler with server handler.
func Handler(hander http.Handler) Option {
	return func(s *Server) {
		s.Handler = hander
	}
}

// Address with server address.
func Address(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

// StopTimeout with server stopTimeout.
func StopTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.stopTimeout = t
	}
}

// IdleTimeout with server IdleTimeout.
func IdleTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.IdleTimeout = t
	}
}

// ReadTimeout with server ReadTimeout.
func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.ReadTimeout = t
	}
}

// WriteTimeout with server WriteTimeout.
func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.WriteTimeout = t
	}
}

type Server struct {
	http.Server
	stopTimeout time.Duration
}

func New(opts ...Option) *Server {
	srv := &Server{
		stopTimeout: 3 * time.Second,
	}

	for _, o := range opts {
		o(srv)
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	log.Println("[HTTP] Server listen:" + s.Addr)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP listen error: %s", err)
		}
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			_ = s.Stop(ctx)
			return ctx.Err()
		case <-c:
			log.Printf("[HTTP] Shutdown waiting %s\n", s.stopTimeout.String())
			time.Sleep(s.stopTimeout)
			err := s.Stop(ctx)
			if err != nil {
				return fmt.Errorf("HTTP Server Shutdown:%s", err)
			}
			return nil
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("[HTTP] server stopping")
	_ = s.Close()
	return s.Shutdown(ctx)
}
