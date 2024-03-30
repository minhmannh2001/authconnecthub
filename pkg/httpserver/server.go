package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.server.ListenAndServe()
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
