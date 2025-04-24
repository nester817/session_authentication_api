package server

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
}

func NewServer(addr string, h http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:    addr,
			Handler: h,

			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,

			IdleTimeout: 120 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
