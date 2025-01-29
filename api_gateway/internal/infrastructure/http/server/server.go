package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

// NewServer creates a new server instance
func NewServer(port string, handler http.Handler, logger *log.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		logger: logger,
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	s.logger.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
